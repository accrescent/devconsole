package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/data"
	"github.com/accrescent/devconsole/quality"
)

func NewUpdate(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	storage := c.MustGet("storage").(data.FileStorage)
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

	// We've received the (supposed) APK set. Now extract the app metadata.
	apk, appFile, err := openAPKSet(file)
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

	m := apk.Manifest()

	// Run tests whose failures warrant immediate rejection
	if m.Package != appID {
		msg := fmt.Sprintf(
			"App ID '%s' doesn't match expected value '%s'",
			appID,
			m.Package,
		)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
		return
	}
	if m.VersionCode <= versionCode {
		err := fmt.Sprintf(
			"Version %d is not more than current version %d",
			m.VersionCode, versionCode,
		)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}
	if err := quality.RunRejectTests(apk, quality.Update); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// App passed all automated checks, so save it to disk
	apkSetHandle, err := storage.SaveUpdate(appFile)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Run tests whose failures warrant manual review
	issues := quality.RunReviewTests(apk)

	if err := db.CreateUpdate(
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
