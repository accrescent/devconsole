package page

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v45/github"

	"github.com/accrescent/devportal/config"
)

func Dashboard(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	client := c.MustGet("gh_client").(*github.Client)
	conf := c.MustGet("config").(*config.Config)

	user, _, err := client.Users.Get(c, "")
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if *user.ID == conf.SignerGitHubID {
		rows, err := db.Query("SELECT id FROM approved_apps")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()

		var appIDs []string
		for rows.Next() {
			var appID string
			if err := rows.Scan(&appID); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			appIDs = append(appIDs, appID)
		}

		c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
			"username":     user.Login,
			"pending_apps": appIDs,
		})
	} else {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"username": user.Login,
		})
	}
}
