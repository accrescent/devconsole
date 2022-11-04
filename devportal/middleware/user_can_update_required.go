package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devportal/data"
)

func UserCanUpdateRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(data.DB)
		ghID := c.MustGet("gh_id").(int64)
		appID := c.Param("id")

		userCanUpdate, err := db.GetUserPermissions(appID, ghID)
		if err != nil {
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
