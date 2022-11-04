package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"

	"github.com/accrescent/devportal/data"
	"github.com/accrescent/devportal/quality"
)

func PublishApp(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	appID := c.Param("id")

	ghID, label, versionCode, versionName, iconID, path, err := db.GetSubmittedAppInfo(appID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Publish to repository server
	if err := publish(c, appID, versionCode, versionName, quality.NewApp, path); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := db.PublishApp(appID, label, versionCode, versionName, iconID, ghID); err != nil {
		if errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
			_ = c.AbortWithError(http.StatusConflict, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.String(http.StatusOK, "")
}
