package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
)

func PublishApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)
	appID := c.Param("appID")

	var ghID int
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

	if _, err := db.Exec("INSERT INTO app_teams (id) VALUES (?)", appID); err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			_ = c.AbortWithError(http.StatusConflict, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	if _, err := db.Exec(
		"INSERT INTO app_team_users (app_id, user_gh_id) VALUES (?, ?)",
		appID, ghID,
	); err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			_ = c.AbortWithError(http.StatusConflict, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	if _, err := db.Exec("DELETE FROM approved_apps WHERE id = ?", appID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
