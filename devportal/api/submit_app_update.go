package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/quality"
)

func SubmitAppUpdate(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)
	stagingAppID, err := c.Cookie(stagingUpdateAppIDCookie)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var label, appPath, versionName string
	var versionCode int
	if err := db.QueryRow(
		`SELECT label, version_code, version_name, path
		FROM staging_app_updates
		WHERE id = ? AND session_id = ?`,
		stagingAppID, sessionID,
	).Scan(&label, &versionCode, &versionName, &appPath); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	rows, err := db.Query(`
		SELECT review_error_id
		FROM staging_update_review_errors
		WHERE staging_app_id = ?
	`, stagingAppID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	var reviewErrors []string
	for rows.Next() {
		var reviewError string
		if err := rows.Scan(&reviewError); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		reviewErrors = append(reviewErrors, reviewError)
	}

	tx, err := db.Begin()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if len(reviewErrors) > 0 {
		if _, err := tx.Exec(
			`REPLACE INTO submitted_updates (id, label, version_code, version_name, path)
			VALUES (?, ?, ?, ?, ?)`,
			stagingAppID, label, versionCode, versionName, appPath,
		); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			if err := tx.Rollback(); err != nil {
				_ = c.Error(err)
			}
			return
		}

		insertQuery := `INSERT INTO submitted_update_review_errors
			(submitted_app_id, review_error_id) VALUES `
		var inserts []string
		var params []interface{}
		for _, rError := range reviewErrors {
			inserts = append(inserts, "(?, ?)")
			params = append(params, stagingAppID, rError)
		}
		insertQuery = insertQuery + strings.Join(inserts, ",")
		if _, err := tx.Exec(insertQuery, params...); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			if err := tx.Rollback(); err != nil {
				_ = c.Error(err)
			}
			return
		}
	} else {
		// No review necessary, so publish the update immediately.
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

		if err := publish(c, stagingAppID, versionCode, versionName,
			quality.AppUpdate, appPath,
		); err != nil {
			if err := tx.Rollback(); err != nil {
				_ = c.Error(err)
			}
			return
		}
	}
	if _, err := tx.Exec(
		"DELETE FROM staging_app_updates WHERE id = ? AND session_id = ?",
		stagingAppID, sessionID,
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
