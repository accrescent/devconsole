package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApproveApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	appID := c.Param("id")

	if _, err := db.Exec(
		"UPDATE submitted_apps SET approved = TRUE WHERE id = ?",
		appID,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
