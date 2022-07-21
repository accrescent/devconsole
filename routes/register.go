package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)

	var input struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		return
	}

	// Verify user is allowed to register with submitted email
	var valid bool
	if err := db.QueryRow(
		`SELECT EXISTS (SELECT 1
		FROM usable_email_cache
		WHERE session_id = ?
		AND email = ?
	)`, sessionID, input.Email).Scan(&valid); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if valid {
		if _, err := db.Exec(
			"DELETE FROM usable_email_cache WHERE session_id = ?",
			sessionID, input.Email,
		); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		_ = c.AbortWithError(http.StatusForbidden, errors.New("email not usable"))
		return
	}

	// Register user
	var ghID string
	if err := db.QueryRow(
		"SELECT gh_id FROM sessions WHERE id = ?",
		sessionID,
	).Scan(&ghID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res, err := db.Exec(
		"INSERT INTO users (gh_id, email) VALUES (?, ?)",
		ghID, input.Email,
	)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if rows != 1 {
		_ = c.AbortWithError(http.StatusConflict, errors.New("user already exists"))
		return
	}

	c.String(http.StatusOK, "")
}
