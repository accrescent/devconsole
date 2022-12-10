package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/auth"
	"github.com/accrescent/devconsole/config"
)

func SignerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := c.MustGet("config").(config.Config)
		ghID := c.MustGet("gh_id").(int64)

		if auth.ConstantTimeEqInt64(ghID, conf.SignerGitHubID) == 0 {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
