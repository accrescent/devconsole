package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

func GitHub(c *gin.Context) {
	conf := c.MustGet("oauth2_config").(*oauth2.Config)

	// CSRF protection. http://tools.ietf.org/html/rfc6749#section-10.12
	state := make([]byte, 16)
	if _, err := rand.Read(state); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	stateStr := hex.EncodeToString(state)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(authStateCookie, stateStr, 30, "/", "", true, true)

	url := conf.AuthCodeURL(stateStr)
	c.Redirect(http.StatusFound, url)
}

func GitHubCallback(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	conf := c.MustGet("oauth2_config").(*oauth2.Config)

	stateParam, exists := c.GetQuery("state")
	if !exists {
		c.SetCookie(authStateCookie, "", -1, "/", "", true, true)
		_ = c.AbortWithError(http.StatusBadRequest, ErrNoStateParam)
		return
	}

	stateCookie, err := c.Cookie(authStateCookie)
	if err != nil {
		if err != http.ErrNoCookie {
			c.SetCookie(authStateCookie, "", -1, "/", "", true, true)
		}
		_ = c.AbortWithError(http.StatusForbidden, err)
		return
	}

	// CSRF protection. http://tools.ietf.org/html/rfc6749#section-10.12
	if stateParam != stateCookie {
		c.SetCookie(authStateCookie, "", -1, "/", "", true, true)
		_ = c.AbortWithError(http.StatusForbidden, ErrNoStateMatch)
		return
	}

	code := c.Query("code")
	token, err := conf.Exchange(c, code)
	if err != nil {
		var retrieveError *oauth2.RetrieveError
		if errors.As(err, &retrieveError) {
			_ = c.AbortWithError(retrieveError.Response.StatusCode, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	// Get authenticated user
	httpClient := conf.Client(c, token)
	client := github.NewClient(httpClient)
	user, _, err := client.Users.Get(c, "")
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// Add session
	sid := make([]byte, 16)
	if _, err := rand.Read(sid); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	sidStr := hex.EncodeToString(sid)

	now := time.Now()
	if _, err := db.Exec("DELETE FROM sessions WHERE expiry_time < ?", now.Unix()); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, err := db.Exec(
		"INSERT INTO sessions (id, gh_id, access_token, expiry_time) VALUES (?, ?, ?, ?)",
		sidStr, user.ID, token.AccessToken, now.Add(24*time.Hour).Unix(),
	); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(SessionCookie, sidStr, 24*60*60, "/", "", true, true) // Max-Age 1 day

	var registered bool
	if err = db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM users WHERE gh_id = ?)",
		user.ID,
	).Scan(&registered); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if registered {
		c.Redirect(http.StatusFound, "/auth/redirect/dashboard")
	} else {
		c.Redirect(http.StatusFound, "/auth/redirect/register")
	}
}
