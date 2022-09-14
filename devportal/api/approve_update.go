package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/quality"
)

func ApproveUpdate(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)
	appID := c.Param("id")

	// Check for authorization
	var reviewer bool
	if err := db.QueryRow(
		"SELECT reviewer FROM users WHERE gh_id = ?",
		ghID,
	).Scan(&reviewer); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !reviewer {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var versionCode int
	var versionName, path string
	if err := db.QueryRow(
		"SELECT version_code, version_name, path FROM submitted_updates WHERE id = ?",
		appID,
	).Scan(&versionCode, &versionName, &path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	tx, err := db.Begin()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err := tx.Exec(
		"UPDATE app_teams SET version_code = ?, version_name = ?",
		versionCode, versionName,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		if err := tx.Rollback(); err != nil {
			_ = c.Error(err)
		}
		return
	}
	if _, err := tx.Exec("DELETE FROM submitted_updates WHERE id = ?", appID); err != nil {
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
