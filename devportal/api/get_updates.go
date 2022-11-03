package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func GetUpdates(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)

	dbApps, err := db.Query(
		`SELECT app_id, label, version_code, version_name, issue_group_id
		FROM submitted_updates
		WHERE reviewer_gh_id = ?`,
		ghID,
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
		var issueGroupID *int
		if err := dbApps.Scan(
			&appID,
			&label,
			&versionCode,
			&versionName,
			&issueGroupID,
		); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		dbIssues, err := db.Query(
			"SELECT id FROM issues WHERE issue_group_id = ?",
			issueGroupID,
		)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer dbIssues.Close()
		var issues []string
		for dbIssues.Next() {
			var issueID string
			if err := dbIssues.Scan(&issueID); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			issues = append(issues, issueID)
		}

		app := data.App{
			AppID:       appID,
			Label:       label,
			VersionCode: versionCode,
			VersionName: versionName,
			Issues:      issues,
		}
		apps = append(apps, app)
	}

	c.JSON(http.StatusOK, apps)
}
