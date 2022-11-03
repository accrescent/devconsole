package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"

	"github.com/accrescent/devportal/quality"
)

func SubmitUpdate(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)
	appID := c.Param("id")
	versionCode, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var label, appPath, versionName string
	var issueGroupID *int
	if err := db.QueryRow(
		`SELECT label, version_name, path, issue_group_id
		FROM staging_updates
		WHERE app_id = ? AND version_code = ? AND user_gh_id = ?`,
		appID, versionCode, ghID,
	).Scan(&label, &versionName, &appPath, &issueGroupID); err != nil {
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
	if issueGroupID != nil {
		if _, err := tx.Exec(
			`INSERT INTO submitted_updates (
				app_id,
				label,
				version_code,
				version_name,
				reviewer_gh_id,
				path,
				issue_group_id
			)
			VALUES (
				?,
				?,
				?,
				?,
				(SELECT user_gh_id FROM reviewers ORDER BY RANDOM() LIMIT 1),
				?,
				?
			)`,
			appID, label, versionCode, versionName, appPath, issueGroupID,
		); err != nil {
			if errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintUnique) {
				_ = c.AbortWithError(http.StatusConflict, err)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			if err := tx.Rollback(); err != nil {
				_ = c.Error(err)
			}
			return
		}
	} else {
		// No review necessary, so publish the update immediately.
		if err := publish(c, appID, versionCode, versionName,
			quality.Update, appPath,
		); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if _, err := tx.Exec(
			`UPDATE published_apps
			SET version_code = ?, version_name = ?
			WHERE id = ?`,
			versionCode, versionName,
			appID,
		); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			if err := tx.Rollback(); err != nil {
				_ = c.Error(err)
			}
			return
		}
	}
	if _, err := tx.Exec(
		"DELETE FROM staging_updates WHERE app_id = ? AND version_code = ?",
		appID, versionCode,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
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
