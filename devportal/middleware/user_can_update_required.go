package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserCanUpdateRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)
		ghID := c.MustGet("gh_id").(int64)
		appID := c.Param("id")

		var userCanUpdate bool
		if err := db.QueryRow(
			`SELECT can_update FROM user_permissions
			WHERE app_id = ? AND user_gh_id = ?`,
			appID,
			ghID,
		).Scan(&userCanUpdate); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_ = c.AbortWithError(http.StatusForbidden, err)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}
		if !userCanUpdate {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
