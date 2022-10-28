package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
)

func SubmitApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)
	appID := c.Param("id")

	var label, path, versionName string
	var versionCode int
	if err := db.QueryRow(
		`SELECT label, version_code, version_name, path
		FROM staging_apps
		WHERE id = ? AND user_gh_id = ?`,
		appID, ghID,
	).Scan(&label, &versionCode, &versionName, &path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	rows, err := db.Query(`
		SELECT review_error_id
		FROM staging_app_review_errors
		WHERE staging_app_id = ?
	`, appID)
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
	if _, err := tx.Exec(
		`INSERT INTO submitted_apps (
			id,
			gh_id,
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
			?,
			(SELECT user_gh_id FROM reviewers ORDER BY RANDOM() LIMIT 1),
			?
		)`,
		appID, ghID, label, versionCode, versionName, path,
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
	if len(reviewErrors) > 0 {
		insertQuery := `INSERT INTO submitted_app_review_errors
			(submitted_app_id, review_error_id) VALUES `
		var inserts []string
		var params []interface{}
		for _, rError := range reviewErrors {
			inserts = append(inserts, "(?, ?)")
			params = append(params, appID, rError)
		}
		insertQuery = insertQuery + strings.Join(inserts, ",")
		if _, err := tx.Exec(insertQuery, params...); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			if err := tx.Rollback(); err != nil {
				_ = c.Error(err)
			}
			return
		}
	}
	if _, err := tx.Exec(
		"DELETE FROM staging_apps WHERE id = ? AND user_gh_id = ?",
		appID, ghID,
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
