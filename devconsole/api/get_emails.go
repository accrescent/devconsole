package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v48/github"

	"github.com/accrescent/devconsole/data"
)

func GetEmails(c *gin.Context) {
	ghClient := c.MustGet("gh_client").(*github.Client)

	usableEmails, err := data.GetUsableEmails(c, ghClient)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, usableEmails)
}
