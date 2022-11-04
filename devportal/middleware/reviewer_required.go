package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func ReviewerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(data.DB)
		ghID := c.MustGet("gh_id").(int64)

		_, reviewer, err := db.GetUserRoles(ghID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !reviewer {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
