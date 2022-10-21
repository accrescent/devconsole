package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/auth"
)

func LogOut(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	sessionID := c.MustGet("session_id").(string)

	if _, err := db.Exec("DELETE FROM sessions WHERE id = ?", sessionID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(auth.SessionCookie, "", -1, "/", "", true, true)

	c.String(http.StatusOK, "")
}
