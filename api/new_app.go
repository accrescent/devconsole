package api

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"errors"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/accrescent/apkstat"
	"github.com/gin-gonic/gin"
)

func NewApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)

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

	preParsedDir := filepath.Join("pre_parsed", ghID)
	if err := os.MkdirAll(preParsedDir, 0500); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	preParsedDir, err := os.MkdirTemp(preParsedDir, "")
	if err != nil && !errors.Is(err, fs.ErrExist) {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// file.Filename shouldn't be trusted, so we strip all path separators before saving. See
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition#directives
	// and https://github.com/gin-gonic/gin/issues/1693.
	filename := filepath.Base(file.Filename)
	upload := filepath.Join(preParsedDir, filename)
	defer os.Remove(upload)
	if err := c.SaveUploadedFile(file, upload); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// We've received the (supposed) APK set. Now extract the app metadata.
	apkSet, err := zip.OpenReader(upload)
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

	uploadKey := strconv.FormatUint(rand.Uint64(), 16)
	validDir := filepath.Join("valid", ghID, uploadKey)
	if err := os.MkdirAll(validDir, 0500); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	validPath := filepath.Join(validDir, filename)
	if err := os.Rename(upload, validPath); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := os.RemoveAll(preParsedDir); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	m := apk.Manifest()

	if _, err := db.Exec(
		`INSERT INTO valid_apps (
			gh_id, upload_key, path, id,
			label, version_code, version_name
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		ghID, uploadKey, validPath, m.Package,
		m.Application.Label, m.VersionCode, m.VersionName,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("__Host-upload_key", uploadKey, 5*60, "/", "", true, true) // Max-Age 5 min

	c.JSON(http.StatusCreated, gin.H{
		"id":           m.Package,
		"label":        m.Application.Label,
		"version_name": m.VersionName,
		"version_code": m.VersionCode,
	})
}
