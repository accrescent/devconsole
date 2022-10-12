package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReviewerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*sql.DB)
		ghID := c.MustGet("gh_id").(int64)

		var reviewer bool
		if err := db.QueryRow(
			"SELECT EXISTS (SELECT 1 from reviewers WHERE user_gh_id = ?)",
			ghID,
		).Scan(&reviewer); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !reviewer {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
