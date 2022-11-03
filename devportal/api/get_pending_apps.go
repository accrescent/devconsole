package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	AppID       string   `json:"app_id"`
	Label       string   `json:"label"`
	VersionCode int      `json:"version_code"`
	VersionName string   `json:"version_name"`
	Issues      []string `json:"issues,omitempty"`
}

func GetPendingApps(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)

	dbApps, err := db.Query(
		`SELECT id, label, version_code, version_name, issue_group_id
		FROM submitted_apps
		WHERE reviewer_gh_id = ? AND approved = FALSE`,
		ghID,
	)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer dbApps.Close()
	var apps []App
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

		app := App{appID, label, versionCode, versionName, issues}
		apps = append(apps, app)
	}

	c.JSON(http.StatusOK, apps)
}
