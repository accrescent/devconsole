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
	appID := c.Param("id")

	app, _, _, _, handle, err := db.GetSubmittedAppInfo(appID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Publish to repository server
	if err := publish(c, appID, app.VersionCode, app.VersionName, quality.NewApp, handle); err != nil {
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
