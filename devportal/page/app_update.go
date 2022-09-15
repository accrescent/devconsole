package page

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppUpdate(c *gin.Context) {
	appID := c.Param("id")

	c.HTML(http.StatusOK, "app_update.html", gin.H{
		"id": appID,
	})
}
