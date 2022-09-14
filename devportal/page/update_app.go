package page

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateApp(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	appID := c.Param("id")

	var label, versionName string
	var versionCode int
	if err := db.QueryRow(
		`SELECT label, version_code, version_name
		FROM app_teams WHERE id = ?`,
		appID,
	).Scan(&label, &versionCode, &versionName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.HTML(http.StatusOK, "update_app.html", gin.H{
		"id":           appID,
		"label":        label,
		"version_code": versionCode,
		"version_name": versionName,
	})
}
