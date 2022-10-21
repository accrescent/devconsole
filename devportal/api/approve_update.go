package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/quality"
)

func ApproveUpdate(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	appID := c.Param("id")
	version := c.Param("version")
	versionCode, err := strconv.Atoi(version)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var firstUpdateVersion int
	var versionName, path string
	if err := db.QueryRow(
		`SELECT (SELECT MIN(version_code) FROM submitted_updates), version_name, path
		FROM submitted_updates
		WHERE app_id = ? AND version_code = ?`,
		appID,
		versionCode,
	).Scan(&firstUpdateVersion, &versionName, &path); err != nil {
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

	tx, err := db.Begin()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err := tx.Exec(
		"UPDATE published_apps SET version_code = ?, version_name = ?",
		versionCode, versionName,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		if err := tx.Rollback(); err != nil {
			_ = c.Error(err)
		}
		return
	}
	if _, err := tx.Exec(
		"DELETE FROM submitted_updates WHERE app_id = ? AND version_code = ?",
		appID,
		versionCode,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		if err := tx.Rollback(); err != nil {
			_ = c.Error(err)
		}
		return
	}

	if err := publish(c, appID, versionCode, versionName, quality.AppUpdate, path); err != nil {
		if err := tx.Rollback(); err != nil {
			_ = c.Error(err)
		}
		return
	}

	if err := tx.Commit(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
