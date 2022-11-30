package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func GetUpdateAPKs(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	storage := c.MustGet("storage").(data.FileStorage)
	appID := c.Param("id")
	versionCodeStr := c.Param("version")
	versionCode, err := strconv.Atoi(versionCodeStr)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, _, handle, err := db.GetUpdateInfo(appID, versionCode)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	file, size, err := storage.GetAPKSet(handle)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	filename := appID + "-" + versionCodeStr + ".apks"
	headers := map[string]string{"Content-Disposition": `attachment; filename="` + filename + `"`}

	c.DataFromReader(http.StatusOK, size, "application/octet-stream", file, headers)
}
