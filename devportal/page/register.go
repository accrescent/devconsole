package page

import (
	"database/sql"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v47/github"
)

// GitHub's noreply email address format documented at
// https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/setting-your-commit-email-address
var noReplyEmail = regexp.MustCompile(`^([0-9]{7}\+)?.*@users\.noreply\.github\.com$`)

func Register(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)
	client := c.MustGet("gh_client").(*github.Client)

	emails, _, err := client.Users.ListEmails(c, nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Filter out unverified and noreply emails
	var usableEmails []string
	for _, email := range emails {
		address := email.GetEmail()
		if email.GetVerified() && !noReplyEmail.MatchString(address) {
			usableEmails = append(usableEmails, email.GetEmail())
		}
	}

	// Cache usable emails for later validation
	emailInsertQuery := "INSERT OR IGNORE INTO usable_email_cache (session_id, email) VALUES "
	var inserts []string
	var params []interface{}
	for _, email := range usableEmails {
		inserts = append(inserts, "(?, ?)")
		params = append(params, sessionID, email)
	}
	emailInsertQuery = emailInsertQuery + strings.Join(inserts, ",")
	if _, err := db.Exec(emailInsertQuery, params...); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "register.html", gin.H{
		"emails": usableEmails,
	})
}
