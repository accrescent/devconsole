package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func RejectApp(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	appID := c.Param("id")

	if err := db.DeleteSubmittedApp(appID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
