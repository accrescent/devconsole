package api

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"image"
	_ "image/png"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/data"
	"github.com/accrescent/devconsole/quality"
)

func NewApp(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	storage := c.MustGet("storage").(data.FileStorage)
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

	// We've received the (supposed) APK set. Now extract the app metadata.
	metadata, apk, appFile, err := openAPKSet(formApp)
	if err != nil {
		if errors.Is(err, ErrFatalIO) {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			msg := "App is in incorrect format. Make sure you upload an APK set."
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
		}
		return
	}
	defer appFile.Close()

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

	// Run tests whose failures warrant immediate rejection
	if err := quality.RunRejectTests(metadata, apk, quality.NewApp); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	m := apk.Manifest()

	// If app already exists on disk, delete it
	overwrite := true
	appHandle, iconHandle, err := db.GetStagingAppInfo(m.Package, ghID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			overwrite = false
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	if overwrite {
		if err := storage.DeleteApp(appHandle); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if err := storage.DeleteIcon(iconHandle); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// App passed all automated checks, so save it to disk
	if _, err := formIconFile.Seek(0, io.SeekStart); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	apkSetHandle, iconHandle, err := storage.SaveNewApp(appFile, formIconFile)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Run tests whose failures warrant manual review
	issues := quality.RunReviewTests(apk)

	// Calculate icon hash
	if _, err := formIconFile.Seek(0, io.SeekStart); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	hasher := sha256.New()
	if _, err := io.Copy(hasher, formIconFile); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	iconHash := hex.EncodeToString(hasher.Sum(nil))

	if err := db.CreateApp(
		data.AppWithIssues{
			App: data.App{
				AppID:       m.Package,
				Label:       *m.Application.Label,
				VersionCode: m.VersionCode,
				VersionName: m.VersionName,
			},
			Issues: issues,
		},
		ghID,
		apkSetHandle,
		iconHandle,
		iconHash,
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
