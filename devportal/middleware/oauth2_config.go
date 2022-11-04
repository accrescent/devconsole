package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func OAuth2Config(conf oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("oauth2_config", conf)
		c.Next()
	}
}
