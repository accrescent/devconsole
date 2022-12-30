package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/data"
	"github.com/accrescent/devconsole/quality"
)

func ApproveUpdate(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	storage := c.MustGet("storage").(data.FileStorage)
	appID := c.Param("id")
	versionCode, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	firstUpdate, versionName, appHandle, issueGroupID, err := db.GetUpdateInfo(appID, versionCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	// Prohibit approving updates out-of-order
	if versionCode != firstUpdate {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	if err := publish(
		c,
		appID,
		int32(versionCode),
		versionName,
		quality.Update,
		appHandle,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete local copy of update once it's published
	if err := storage.DeleteApp(appHandle); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := db.ApproveUpdate(appID, versionCode, versionName, issueGroupID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
