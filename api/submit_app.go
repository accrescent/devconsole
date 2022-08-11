package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
)

func SubmitApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)
	ghID := c.MustGet("gh_id").(int64)
	stagingAppID, err := c.Cookie(stagingAppIDCookie)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var path string
	if err := db.QueryRow(
		"SELECT path FROM staging_apps WHERE id = ? AND session_id = ?",
		stagingAppID, sessionID,
	).Scan(&path); err != nil {
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
	if _, err := tx.Exec(
		"INSERT INTO approved_apps (id, gh_id, path) VALUES (?, ?, ?)",
		stagingAppID, ghID, path,
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
	if _, err := tx.Exec(
		"DELETE FROM staging_apps WHERE id = ? AND session_id = ?",
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
