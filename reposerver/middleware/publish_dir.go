package middleware

import "github.com/gin-gonic/gin"

func PublishDir(publishDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("publish_dir", publishDir)
		c.Next()
	}
}
