package api

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/accrescent/apkstat"
	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/quality"
)

func NewApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)

	file, err := c.FormFile("file")
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
			"SELECT EXISTS (SELECT 1 FROM staging_apps WHERE session_id = ? AND path = ?)",
			sessionID, filename,
		).Scan(&unsubmitted); err != nil {
			_ = cCp.Error(err)
			return
		}

		if unsubmitted {
			if _, err := db.Exec(
				"DELETE FROM staging_apps WHERE session_id = ? AND path = ?",
				sessionID, filename,
			); err != nil {
				_ = cCp.Error(err)
			}
			os.RemoveAll(dir)
		}
	}()

	// We've received the (supposed) APK set. Now extract the app metadata.
	apkSet, err := zip.OpenReader(filename)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer apkSet.Close()
	rawBaseAPK, err := apkSet.Open("splits/base-master.apk")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer rawBaseAPK.Close()
	baseAPK, err := io.ReadAll(rawBaseAPK)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	apk, err := apk.FromReader(bytes.NewReader(baseAPK), int64(len(baseAPK)))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Run tests whose failures warrant immediate rejection
	if err := quality.RunRejectTests(apk, quality.NewApp); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Run tests whose failures warrant manual review
	insertQuery := "INSERT OR IGNORE INTO review_errors (id) VALUES "
	var inserts []string
	var params []interface{}
	reviewErrors := quality.RunReviewTests(apk)
	for _, rError := range reviewErrors {
		inserts = append(inserts, "(?)")
		params = append(params, rError)
	}
	insertQuery = insertQuery + strings.Join(inserts, ",")
	if len(reviewErrors) > 0 {
		if _, err := db.Exec(insertQuery, params...); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	m := apk.Manifest()

	tx, err := db.Begin()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, err := tx.Exec(
		"REPLACE INTO staging_apps (id, session_id, path) VALUES (?, ?, ?)",
		m.Package, sessionID, filename,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		if err := tx.Rollback(); err != nil {
			_ = c.Error(err)
		}
		return
	}
	if len(reviewErrors) > 0 {
		insertQuery = `INSERT INTO staging_app_review_errors
		(staging_app_id, staging_app_session_id, review_error_id) VALUES `
		inserts = []string{}
		params = []interface{}{}
		for _, rError := range reviewErrors {
			inserts = append(inserts, "(?, ?, ?)")
			params = append(params, m.Package, sessionID, rError)
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
	if err := tx.Commit(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(stagingAppIDCookie, m.Package, 5*60, "/", "", true, true) // Max-Age 5 min

	c.JSON(http.StatusCreated, gin.H{
		"id":            m.Package,
		"label":         m.Application.Label,
		"version_name":  m.VersionName,
		"version_code":  m.VersionCode,
		"review_errors": reviewErrors,
	})
}
