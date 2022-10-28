package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"

	"github.com/accrescent/devportal/quality"
)

func SubmitAppUpdate(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)
	stagingUpdateID, err := c.Cookie(stagingUpdateIDCookie)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var appID, label, appPath, versionName string
	var versionCode int
	if err := db.QueryRow(
		`SELECT app_id, label, version_code, version_name, path
		FROM staging_app_updates
		WHERE id = ? AND user_gh_id = ?`,
		stagingUpdateID, ghID,
	).Scan(&appID, &label, &versionCode, &versionName, &appPath); err != nil {
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
	`, stagingUpdateID)
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
			`INSERT INTO submitted_updates (
				app_id,
				label,
				version_code,
				version_name,
				reviewer_gh_id,
				path
			)
			VALUES (
				?,
				?,
				?,
				?,
				(SELECT user_gh_id FROM reviewers ORDER BY RANDOM() LIMIT 1),
				?
			)`,
			appID, label, versionCode, versionName, appPath,
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

		insertQuery := `INSERT INTO submitted_update_review_errors
			(submitted_app_id, review_error_id) VALUES `
		var inserts []string
		var params []interface{}
		for _, rError := range reviewErrors {
			inserts = append(inserts, "(?, ?)")
			params = append(params, stagingUpdateID, rError)
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
		if err := publish(c, appID, versionCode, versionName,
			quality.AppUpdate, appPath,
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
		"DELETE FROM staging_app_updates WHERE app_id = ? AND version_code = ?",
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
