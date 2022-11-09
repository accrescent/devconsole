package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func GetAppAPKs(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	appID := c.Param("id")

	_, _, _, _, _, path, err := db.GetSubmittedAppInfo(appID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.FileAttachment(path, appID+".apks")
}
