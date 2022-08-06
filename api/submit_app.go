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
	stagingAppID, err := c.Cookie(stagingAppIDCookie)
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
	if _, err := db.Exec(
		"INSERT INTO approved_apps (id, gh_id, path) VALUES (?, ?, ?)",
		stagingAppID, ghID, path,
	); err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			_ = c.AbortWithError(http.StatusConflict, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	if _, err := db.Exec(
		"DELETE FROM staging_apps WHERE id = ? AND session_id = ?",
		stagingAppID, sessionID,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
