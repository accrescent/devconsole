package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RejectApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	appID := c.Param("id")

	if _, err := db.Exec("DELETE FROM submitted_apps WHERE id = ?", appID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
