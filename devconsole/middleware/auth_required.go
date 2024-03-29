package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"

	"github.com/accrescent/devconsole/auth"
	"github.com/accrescent/devconsole/data"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(data.DB)
		conf := c.MustGet("oauth2_config").(oauth2.Config)

		sessionID, err := c.Cookie(auth.SessionCookie)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		ghID, token, err := db.GetSessionInfo(sessionID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_ = c.AbortWithError(http.StatusUnauthorized, err)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		httpClient := conf.Client(c, &oauth2.Token{AccessToken: token})
		client := github.NewClient(httpClient)

		c.Set("session_id", sessionID)
		c.Set("gh_id", ghID)
		c.Set("gh_client", client)

		c.Next()
	}
}
