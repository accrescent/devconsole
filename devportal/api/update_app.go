package api

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
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

func UpdateApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)
	appID := c.Param("id")

	var versionCode int
	var versionName string
	if err := db.QueryRow(
		"SELECT version_code, version_name from app_teams WHERE id = ?",
		appID,
	).Scan(&versionCode, &versionName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

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
			"SELECT EXISTS (SELECT 1 FROM staging_app_updates WHERE id = ?)",
			appID,
		).Scan(&unsubmitted); err != nil {
			_ = cCp.Error(err)
			return
		}

		if unsubmitted {
			if _, err := db.Exec(
				"DELETE FROM staging_app_updates WHERE id = ?",
				appID,
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
	if err := quality.RunRejectTests(apk, quality.AppUpdate); err != nil {
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

	tx, err := db.Begin()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, err := tx.Exec(
		`REPLACE INTO staging_app_updates (id, session_id, label, version_code, version_name, path)
		VALUES (?, ?, ?, ?, ?, ?)`,
		m.Package, sessionID, m.Application.Label, m.VersionCode, m.VersionName, filename,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		if err := tx.Rollback(); err != nil {
			_ = c.Error(err)
		}
		return
	}
	if len(reviewErrors) > 0 {
		insertQuery = `INSERT INTO staging_update_review_errors
		(staging_app_id, staging_update_session_id, review_error_id) VALUES `
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
	c.SetCookie(stagingUpdateAppIDCookie, m.Package, 5*60, "/", "", true, true) // Max-Age 5 min

	c.JSON(http.StatusCreated, gin.H{
		"id":            m.Package,
		"current_vcode": versionCode,
		"current_vname": versionName,
		"new_vcode":     m.VersionCode,
		"new_vname":     m.VersionName,
		"review_errors": reviewErrors,
	})
}
