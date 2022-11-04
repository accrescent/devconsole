package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func GetApprovedApps(c *gin.Context) {
	db := c.MustGet("db").(data.DB)

	apps, err := db.GetApprovedApps()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, apps)
}
