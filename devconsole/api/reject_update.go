package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/data"
)

func RejectUpdate(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	storage := c.MustGet("storage").(data.FileStorage)
	appID := c.Param("id")
	versionCode, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, _, appHandle, _, err := db.GetUpdateInfo(appID, versionCode)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := db.DeleteSubmittedUpdate(appID, versionCode); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := storage.DeleteFile(appHandle); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
