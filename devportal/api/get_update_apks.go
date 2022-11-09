package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func GetUpdateAPKs(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	appID := c.Param("id")
	versionCodeStr := c.Param("version")
	versionCode, err := strconv.Atoi(versionCodeStr)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, _, path, err := db.GetUpdateInfo(appID, versionCode)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.FileAttachment(path, appID+"-"+versionCodeStr+".apks")
}
