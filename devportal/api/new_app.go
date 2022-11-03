package api

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/quality"
)

func NewApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)

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
	// Delete app after 5 minutes if not submitted
	cCp := c.Copy()
	go func() {
		time.Sleep(5 * time.Minute)

		var unsubmitted bool
		if err := db.QueryRow(
			"SELECT EXISTS (SELECT 1 FROM staging_apps WHERE user_gh_id = ? AND path = ?)",
			ghID, filename,
		).Scan(&unsubmitted); err != nil {
			_ = cCp.Error(err)
			return
		}

		if unsubmitted {
			if _, err := db.Exec(
				"DELETE FROM staging_apps WHERE user_gh_id = ? AND path = ?",
				ghID, filename,
			); err != nil {
				_ = cCp.Error(err)
			}
			os.RemoveAll(dir)
		}
	}()

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

	// Run tests whose failures warrant immediate rejection
	if err := quality.RunRejectTests(apk, quality.NewApp); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Run tests whose failures warrant manual review
	issues := quality.RunReviewTests(apk)

	m := apk.Manifest()

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
			params = append(params, issue, groupID)
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
		`REPLACE INTO staging_apps (
			id,
			user_gh_id,
			label,
			version_code,
			version_name,
			path,
			issue_group_id
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

	c.JSON(http.StatusCreated, gin.H{
		"app_id":       m.Package,
		"label":        m.Application.Label,
		"version_name": m.VersionName,
		"version_code": m.VersionCode,
	})
}
