package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/data"
)

func RejectApp(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	storage := c.MustGet("storage").(data.FileStorage)
	appID := c.Param("id")

	_, _, _, _, handle, err := db.GetSubmittedAppInfo(appID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := db.DeleteSubmittedApp(appID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := storage.DeleteApp(handle); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
