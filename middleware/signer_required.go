package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/config"
)

func SignerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)
		conf := c.MustGet("config").(*config.Config)
		sessionID := c.MustGet("session_id").(string)

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

		if ghID != conf.SignerGitHubID {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
