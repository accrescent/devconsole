package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v48/github"
	"golang.org/x/exp/slices"

	"github.com/accrescent/devportal/data"
)

func Register(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	ghID := c.MustGet("gh_id").(int64)
	ghClient := c.MustGet("gh_client").(*github.Client)

	var input struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		return
	}

	// Verify user is allowed to register with the submitted email
	usableEmails, err := data.GetUsableEmails(c, ghClient)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !slices.Contains(usableEmails, input.Email) {
		_ = c.AbortWithError(http.StatusForbidden, errors.New("email not usable"))
		return
	}

	// Register user
	res, err := db.Exec(
		"INSERT INTO users (gh_id, email) VALUES (?, ?)",
		ghID, input.Email,
	)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if rows != 1 {
		_ = c.AbortWithError(http.StatusConflict, errors.New("user already exists"))
		return
	}

	c.String(http.StatusOK, "")
}
