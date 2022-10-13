package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"

	"github.com/accrescent/devportal/quality"
)

func PublishApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)
	appID := c.Param("id")

	var label, appPath, versionName string
	var versionCode int
	if err := db.QueryRow(
		"SELECT label, version_code, version_name, path FROM submitted_apps WHERE id = ?",
		appID,
	).Scan(&label, &versionCode, &versionName, &appPath); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, err := tx.Exec(`INSERT INTO published_apps (id, label, version_code, version_name)
		VALUES (?, ?, ?, ?)`,
		appID, label, versionCode, versionName,
	); err != nil {
		if errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
			_ = c.AbortWithError(http.StatusConflict, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		if err := tx.Rollback(); err != nil {
			_ = c.Error(err)
		}
		return
	}
	if _, err := tx.Exec(
		"INSERT INTO user_permissions (app_id, user_gh_id, can_update) VALUES (?, ?, TRUE)",
		appID, ghID,
	); err != nil {
		if errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
			_ = c.AbortWithError(http.StatusConflict, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		if err := tx.Rollback(); err != nil {
			_ = c.Error(err)
		}
		return
	}
	if _, err := tx.Exec("DELETE FROM submitted_apps WHERE id = ?", appID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		if err := tx.Rollback(); err != nil {
			_ = c.Error(err)
		}
		return
	}

	// Publish to repository server
	if err := publish(c, appID, versionCode, versionName, quality.NewApp, appPath); err != nil {
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
