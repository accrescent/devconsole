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

	isSigner := false
	if *user.ID == conf.SignerGitHubID {
		isSigner = true
	}

	var appIDs []string
	if isSigner {
		// Only display apps not awaiting manual review
		rows, err := db.Query(`SELECT id FROM submitted_apps WHERE NOT EXISTS (
			SELECT 1 FROM submitted_app_review_errors
			WHERE id = submitted_app_id
		)`)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var appID string
			if err := rows.Scan(&appID); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			appIDs = append(appIDs, appID)
		}
	}
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"username":     user.Login,
		"is_signer":    isSigner,
		"pending_apps": appIDs,
	})
}
