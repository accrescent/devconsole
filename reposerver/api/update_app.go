package api

import "github.com/gin-gonic/gin"

func UpdateApp(c *gin.Context) {
	publish(c, appUpdate)
}
