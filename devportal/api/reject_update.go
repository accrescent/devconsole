package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RejectUpdate(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	appID := c.Param("id")
	versionCode, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if _, err := db.Exec(
		"DELETE FROM submitted_updates WHERE app_id = ? AND version_code = ?",
		appID,
		versionCode,
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "")
}
