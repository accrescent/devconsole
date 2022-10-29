package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApproveApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	var json struct {
		AppID string `json:"app_id" binding:"required"`
	}
	if err := c.BindJSON(&json); err != nil {
		return
	}

	if _, err := db.Exec(
		"UPDATE submitted_apps SET approved = TRUE WHERE submitted_app_id = ?",
		json.AppID,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
