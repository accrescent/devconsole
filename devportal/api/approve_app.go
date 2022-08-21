package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApproveApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)

	// Check for authorization
	var reviewer bool
	if err := db.QueryRow(
		"SELECT reviewer FROM users WHERE gh_id = ?",
		ghID,
	).Scan(&reviewer); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !reviewer {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var json struct {
		AppID string `json:"app_id" binding:"required"`
	}
	if err := c.BindJSON(&json); err != nil {
		return
	}

	if _, err := db.Exec(
		"DELETE FROM submitted_app_review_errors WHERE submitted_app_id = ?",
		json.AppID,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
