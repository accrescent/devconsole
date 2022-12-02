package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func GetAppAPKs(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	storage := c.MustGet("storage").(data.FileStorage)
	appID := c.Param("id")

	_, _, _, _, handle, err := db.GetSubmittedAppInfo(appID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	file, size, err := storage.GetAPKSet(handle)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	filename := appID + ".apks"
	headers := map[string]string{"Content-Disposition": `attachment; filename="` + filename + `"`}

	c.DataFromReader(http.StatusOK, size, "application/octet-stream", file, headers)
}
