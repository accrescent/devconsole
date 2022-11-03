package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func GetApps(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)

	rows, err := db.Query(
		`SELECT id, label, version_code, version_name FROM published_apps
		JOIN user_permissions ON user_permissions.app_id = published_apps.id
		WHERE user_permissions.user_gh_id = ?`,
		ghID,
	)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	var apps []data.App
	for rows.Next() {
		var id, label, versionName string
		var versionCode int
		if err := rows.Scan(&id, &label, &versionCode, &versionName); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		app := data.App{
			AppID:       id,
			Label:       label,
			VersionCode: versionCode,
			VersionName: versionName,
			Issues:      []string{},
		}
		apps = append(apps, app)
	}

	c.JSON(http.StatusOK, apps)
}
