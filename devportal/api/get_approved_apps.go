package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func GetApprovedApps(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	dbApps, err := db.Query(
		`SELECT id, label, version_code, version_name
		FROM submitted_apps
		WHERE approved = TRUE`,
	)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer dbApps.Close()
	var apps []data.App
	for dbApps.Next() {
		var appID, label, versionName string
		var versionCode int
		if err := dbApps.Scan(&appID, &label, &versionCode, &versionName); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		app := data.App{
			AppID:       appID,
			Label:       label,
			VersionCode: versionCode,
			VersionName: versionName,
			Issues:      []string{},
		}
		apps = append(apps, app)
	}

	c.JSON(http.StatusOK, apps)
}
