package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/data"
)

func FileStorage(storage data.FileStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("storage", storage)
		c.Next()
	}
}
