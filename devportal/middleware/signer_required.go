package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/config"
)

func SignerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := c.MustGet("config").(*config.Config)
		ghID := c.MustGet("gh_id").(int64)

		if ghID != conf.SignerGitHubID {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
