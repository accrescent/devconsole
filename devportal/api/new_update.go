package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/quality"
)

func NewUpdate(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)
	appID := c.Param("id")

	var versionCode int
	if err := db.QueryRow(
		"SELECT version_code from published_apps WHERE id = ?",
		appID,
	).Scan(&versionCode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	file, err := c.FormFile("app")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	dir, err := os.MkdirTemp("/", "")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	filename := filepath.Join(dir, "app.apks")
	if err := c.SaveUploadedFile(file, filename); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// We've received the (supposed) APK set. Now extract the app metadata.
	apk, err := apkFromAPKSet(filename)
	if err != nil {
		if errors.Is(err, ErrFatalIO) {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	m := apk.Manifest()

	// Run tests whose failures warrant immediate rejection
	if int(m.VersionCode) <= versionCode {
		err := fmt.Sprintf(
			"version %d is not more than current version %d",
			m.VersionCode, versionCode,
		)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}
	if err := quality.RunRejectTests(apk, quality.Update); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Run tests whose failures warrant manual review
	issues := quality.RunReviewTests(apk)

	tx, err := db.Begin()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var issueGroupID *int64
	if len(issues) > 0 {
		res, err := tx.Exec("INSERT INTO issue_groups DEFAULT VALUES")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			if err := tx.Rollback(); err != nil {
				_ = c.Error(err)
			}
			return
		}
		groupID, err := res.LastInsertId()
		issueGroupID = &groupID
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			if err := tx.Rollback(); err != nil {
				_ = c.Error(err)
			}
			return
		}
		insertQuery := "INSERT INTO issues (id, issue_group_id) VALUES "
		var inserts []string
		var params []interface{}
		for _, issue := range issues {
			inserts = append(inserts, "(?, ?)")
			params = append(params, issue, issueGroupID)
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
		`REPLACE INTO staging_updates (
			app_id, user_gh_id, label, version_code, version_name, path, issue_group_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		m.Package,
		ghID,
		m.Application.Label,
		m.VersionCode,
		m.VersionName,
		filename,
		issueGroupID,
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
	// Delete app after 5 minutes if not submitted
	cCp := c.Copy()
	go func() {
		time.Sleep(5 * time.Minute)

		var unsubmitted bool
		if err := db.QueryRow(
			`SELECT EXISTS (
				SELECT 1 FROM staging_updates WHERE app_id = ? AND version_code = ?
			)`,
			appID,
			m.VersionCode,
		).Scan(&unsubmitted); err != nil {
			_ = cCp.Error(err)
			return
		}

		if unsubmitted {
			if _, err := db.Exec(
				"DELETE FROM staging_updates WHERE app_id = ? AND version_code = ?",
				appID,
				m.VersionCode,
			); err != nil {
				_ = cCp.Error(err)
			}
			os.RemoveAll(dir)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"app_id":       m.Package,
		"label":        m.Application.Label,
		"version_code": m.VersionCode,
		"version_name": m.VersionName,
		"issues":       issues,
	})
}
