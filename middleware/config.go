package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/config"
)

func Config(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", conf)
		c.Next()
	}
}
