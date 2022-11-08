package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
	"github.com/accrescent/devportal/quality"
)

func NewApp(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	ghID := c.MustGet("gh_id").(int64)

	formApp, err := c.FormFile("app")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	formIcon, err := c.FormFile("icon")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check that image is a 512x512 PNG
	formIconFile, err := formIcon.Open()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer formIconFile.Close()
	iconInfo, format, err := image.DecodeConfig(formIconFile)
	if err != nil || format != "png" || iconInfo.Width != 512 || iconInfo.Height != 512 {
		msg := "Icon must be a 512x512 PNG"
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	dir, err := os.MkdirTemp("/", "")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	appPath := filepath.Join(dir, "app.apks")
	if err := c.SaveUploadedFile(formApp, appPath); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	iconPath := filepath.Join(dir, "icon.png")
	outIconFile, err := os.Create(iconPath)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer outIconFile.Close()
	if _, err := io.Copy(outIconFile, formIconFile); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// We've received the (supposed) APK set. Now extract the app metadata.
	apk, err := apkFromAPKSet(appPath)
	if err != nil {
		if errors.Is(err, ErrFatalIO) {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			msg := "App is in incorrect format. Make sure you upload an APK set."
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
		}
		return
	}

	// Run tests whose failures warrant immediate rejection
	if err := quality.RunRejectTests(apk, quality.NewApp); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Run tests whose failures warrant manual review
	issues := quality.RunReviewTests(apk)

	m := apk.Manifest()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, outIconFile); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	iconHash := hex.EncodeToString(hasher.Sum(nil))

	if err := db.CreateApp(
		m.Package,
		ghID,
		*m.Application.Label,
		m.VersionCode,
		m.VersionName,
		appPath,
		iconPath,
		iconHash,
		issues,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"app_id":       m.Package,
		"label":        m.Application.Label,
		"version_name": m.VersionName,
		"version_code": m.VersionCode,
	})
}
