package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"

	"github.com/accrescent/devconsole/data"
	"github.com/accrescent/devconsole/quality"
)

func SubmitUpdate(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	ghID := c.MustGet("gh_id").(int64)
	appID := c.Param("id")
	versionCode, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	label, versionName, fileHandle, issueGroupID, needsReview, err := db.GetStagingUpdateInfo(
		appID,
		int32(versionCode),
		ghID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	if !needsReview {
		// No review necessary, so publish the update immediately.
		if err := publish(
			c,
			appID,
			int32(versionCode),
			versionName,
			quality.Update,
			fileHandle,
		); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	app := data.App{
		AppID:       appID,
		Label:       label,
		VersionCode: int32(versionCode),
		VersionName: versionName,
	}
	if err := db.SubmitUpdate(app, fileHandle, issueGroupID, needsReview); err != nil {
		if errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintUnique) {
			_ = c.AbortWithError(http.StatusConflict, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.String(http.StatusOK, "")
}
