package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
	"github.com/accrescent/devportal/quality"
)

func ApproveUpdate(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	appID := c.Param("id")
	version := c.Param("version")
	versionCode, err := strconv.Atoi(version)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	firstUpdateVersion, versionName, path, err := db.GetUpdateInfo(appID, versionCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	// Prohibit approving updates out-of-order
	if versionCode != firstUpdateVersion {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	if err := publish(c, appID, versionCode, versionName, quality.Update, path); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := db.ApproveUpdate(appID, versionCode, versionName); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
