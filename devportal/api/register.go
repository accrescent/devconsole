package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v48/github"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slices"

	"github.com/accrescent/devportal/data"
)

func Register(c *gin.Context) {
	db := c.MustGet("db").(data.DB)
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
	if err := db.CreateUser(ghID, input.Email); err != nil {
		if errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
			_ = c.AbortWithError(http.StatusConflict, err)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.String(http.StatusOK, "")
}
