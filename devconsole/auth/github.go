package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"

	"github.com/accrescent/devconsole/config"
	"github.com/accrescent/devconsole/data"
)

func GitHub(c *gin.Context) {
	conf := c.MustGet("oauth2_config").(oauth2.Config)

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
	db := c.MustGet("db").(data.DB)
	conf := c.MustGet("config").(config.Config)
	oauth2Conf := c.MustGet("oauth2_config").(oauth2.Config)

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
	token, err := oauth2Conf.Exchange(c, code)
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
	httpClient := oauth2Conf.Client(c, token)
	client := github.NewClient(httpClient)
	user, _, err := client.Users.Get(c, "")
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// Check against registration whitelist
	canRegister, err := db.CanUserRegister(*user.ID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if canRegister {
		// Add session
		if err := db.DeleteExpiredSessions(); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		sessionID, err := db.CreateSession(*user.ID, token.AccessToken)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(SessionCookie, sessionID, 24*60*60, "/", "", true, true) // Max-Age 1 day

		registered, reviewer, err := db.GetUserRoles(*user.ID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"logged_in":  true,
			"registered": registered,
			"reviewer":   reviewer,
			"publisher":  *user.ID == conf.SignerGitHubID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"logged_in": false,
		})
	}
}
