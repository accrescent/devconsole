package api

import (
	"database/sql"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SubmitApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)
	uploadKey, err := c.Cookie("__Host-upload_key")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var ghID string
	if err := db.QueryRow(
		"SELECT gh_id FROM sessions WHERE id = ?",
		sessionID,
	).Scan(&ghID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	var appInfo struct {
		path        string
		id          string
		label       string
		versionCode uint32
		versionName string
	}
	if err := db.QueryRow(
		`SELECT path, id, label, version_code, version_name
		FROM valid_apps WHERE gh_id = ? AND upload_key = ?`,
		ghID, uploadKey,
	).Scan(
		&appInfo.path, &appInfo.id, &appInfo.label, &appInfo.versionCode, &appInfo.versionName,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	submittedDir := filepath.Join("submitted", ghID)
	if err := os.MkdirAll(submittedDir, 0500); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	submittedDir, err = os.MkdirTemp(submittedDir, "")
	if err != nil && !errors.Is(err, fs.ErrExist) {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	submittedPath := filepath.Join(submittedDir, filepath.Base(appInfo.path))
	if err := os.Rename(appInfo.path, submittedPath); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := os.RemoveAll(filepath.Dir(appInfo.path)); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, err := db.Exec(
		"DELETE FROM valid_apps WHERE gh_id = ? AND upload_key = ?",
		ghID, uploadKey,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, err := db.Exec(
		`INSERT INTO submitted_apps (
			gh_id, path, id, label, version_code, version_name
		) VALUES (?, ?, ?, ?, ?, ?)`,
		ghID, submittedPath, appInfo.id, appInfo.label, appInfo.versionCode, appInfo.versionName,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
