//go:build !debug

package main

import "github.com/gin-gonic/gin"

func setMode() {
	gin.SetMode(gin.ReleaseMode)
}
