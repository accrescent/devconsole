package api

import (
	"errors"
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

	file, err := c.FormFile("app")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	dir, err := os.MkdirTemp("/", "")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	filename := filepath.Join(dir, "app.apks")
	if err := c.SaveUploadedFile(file, filename); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// We've received the (supposed) APK set. Now extract the app metadata.
	apk, err := apkFromAPKSet(filename)
	if err != nil {
		if errors.Is(err, ErrFatalIO) {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			_ = c.AbortWithError(http.StatusBadRequest, err)
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

	if err := db.CreateApp(
		m.Package,
		ghID,
		*m.Application.Label,
		m.VersionCode,
		m.VersionName,
		filename,
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
