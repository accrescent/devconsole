package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authRequired(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(auth, " ")
		if len(headerParts) != 2 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		authType := headerParts[0]
		if authType != "token" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token := headerParts[1]
		if token != apiKey {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
