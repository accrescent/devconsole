package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"

	"github.com/accrescent/devportal/data"
)

func SubmitApp(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
	ghID := c.MustGet("gh_id").(int64)
	appID := c.Param("id")

	var input struct {
		Label string `json:"label" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		return
	}
	if len(input.Label) < 3 || len(input.Label) > 30 {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if err := db.SubmitApp(appID, input.Label, ghID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := "Nothing to submit. Try uploading and submitting again"
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": msg})
		} else if errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
			msg := "You've already submitted an app with this ID"
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": msg})
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.String(http.StatusOK, "")
}
