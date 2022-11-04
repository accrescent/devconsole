package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func DB(db data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
