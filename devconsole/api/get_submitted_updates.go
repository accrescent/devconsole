package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/data"
)

func GetSubmittedUpdates(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	ghID := c.MustGet("gh_id").(int64)

	updates, err := db.GetSubmittedUpdates(ghID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updates)
}
