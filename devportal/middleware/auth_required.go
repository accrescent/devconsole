package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"

	"github.com/accrescent/devportal/auth"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)
		conf := c.MustGet("oauth2_config").(*oauth2.Config)

		sessionID, err := c.Cookie(auth.SessionCookie)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			c.Redirect(http.StatusFound, "/")
			return
		}

		var token string
		if err := db.QueryRow(
			"SELECT access_token FROM sessions WHERE id = ? AND expiry_time > ?",
			sessionID, time.Now().Unix(),
		).Scan(&token); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.Abort()
				_ = c.Error(err)
				c.Redirect(http.StatusFound, "/")
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		var ghID int64
		if err := db.QueryRow(
			"SELECT gh_id FROM sessions WHERE id = ?",
			sessionID,
		).Scan(&ghID); err != nil {
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
	}
}
