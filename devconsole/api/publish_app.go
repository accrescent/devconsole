package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"

	"github.com/accrescent/devconsole/data"
	"github.com/accrescent/devconsole/quality"
)

func PublishApp(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	storage := c.MustGet("storage").(data.FileStorage)
	appID := c.Param("id")

	app, _, _, _, appHandle, iconHandle, err := db.GetSubmittedAppInfo(appID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Publish to repository server
	if err := publish(
		c,
		appID,
		app.VersionCode,
		app.VersionName,
		quality.NewApp,
		appHandle,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete local copy of app once it's published
	if err := storage.DeleteApp(appHandle); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := storage.DeleteIcon(iconHandle); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := db.PublishApp(appID); err != nil {
		if errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
			_ = c.AbortWithError(http.StatusConflict, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.String(http.StatusOK, "")
}
