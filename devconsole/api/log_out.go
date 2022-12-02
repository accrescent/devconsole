package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/auth"
	"github.com/accrescent/devconsole/data"
)

func LogOut(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	sessionID := c.MustGet("session_id").(string)

	if err := db.DeleteSession(sessionID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(auth.SessionCookie, "", -1, "/", "", true, true)

	c.String(http.StatusOK, "")
}
