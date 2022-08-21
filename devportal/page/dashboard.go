package page

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v46/github"

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
	var sigAppIDs []string
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
			sigAppIDs = append(sigAppIDs, appID)
		}
	}

	var isReviewer bool
	if err := db.QueryRow(
		"SELECT reviewer FROM users WHERE gh_id = ?",
		*user.ID,
	).Scan(&isReviewer); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	reviewApps := make(map[string][]string)
	if isReviewer {
		ids, err := db.Query(`SELECT id FROM submitted_apps WHERE EXISTS (
			SELECT 1 FROM submitted_app_review_errors
			WHERE id = submitted_app_id
		)`)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer ids.Close()

		for ids.Next() {
			var appID string
			if err := ids.Scan(&appID); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			errors, err := db.Query(`SELECT review_error_id
				FROM submitted_app_review_errors
				WHERE submitted_app_id = ?
			`, appID)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			defer errors.Close()
			var rErrors []string
			for errors.Next() {
				var rError string
				if err := errors.Scan(&rError); err != nil {
					_ = c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
				rErrors = append(rErrors, rError)
			}

			reviewApps[appID] = rErrors
		}
	}

	approved, err := db.Query(`SELECT id FROM submitted_apps
		WHERE NOT EXISTS (
			SELECT 1 FROM submitted_app_review_errors
			WHERE id = submitted_app_id
		) AND gh_id = ?`,
		*user.ID,
	)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer approved.Close()
	var approvedApps []string
	for approved.Next() {
		var approvedApp string
		if err := approved.Scan(&approvedApp); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		approvedApps = append(approvedApps, approvedApp)
	}

	published, err := db.Query("SELECT app_id FROM app_team_users WHERE user_gh_id = ?", *user.ID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer published.Close()
	var publishedApps []string
	for published.Next() {
		var publishedApp string
		if err := published.Scan(&publishedApp); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		publishedApps = append(publishedApps, publishedApp)
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"username":            user.Login,
		"is_signer":           isSigner,
		"pending_sig_apps":    sigAppIDs,
		"is_reviewer":         isReviewer,
		"pending_review_apps": reviewApps,
		"approved_apps":       approvedApps,
		"published_apps":      publishedApps,
	})
}
