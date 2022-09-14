package api

import "github.com/gin-gonic/gin"

func PublishApp(c *gin.Context) {
	publish(c, newApp)
}
