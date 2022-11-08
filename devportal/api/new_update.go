package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
	"github.com/accrescent/devportal/quality"
)

func NewUpdate(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	ghID := c.MustGet("gh_id").(int64)
	appID := c.Param("id")

	versionCode, err := db.GetAppInfo(appID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

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
			msg := "app is in incorrect format. Make sure you upload an APK set."
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
		}
		return
	}

	m := apk.Manifest()

	// Run tests whose failures warrant immediate rejection
	if int(m.VersionCode) <= versionCode {
		err := fmt.Sprintf(
			"version %d is not more than current version %d",
			m.VersionCode, versionCode,
		)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}
	if err := quality.RunRejectTests(apk, quality.Update); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Run tests whose failures warrant manual review
	issues := quality.RunReviewTests(apk)

	if err := db.CreateUpdate(
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
		"version_code": m.VersionCode,
		"version_name": m.VersionName,
		"issues":       issues,
	})
}
