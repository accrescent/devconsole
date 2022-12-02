package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/data"
)

func ApproveApp(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	appID := c.Param("id")

	if err := db.ApproveApp(appID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
