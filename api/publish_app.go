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
	ghID := c.MustGet("gh_id").(int)
	appID := c.Param("appID")

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
