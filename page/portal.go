package page

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v45/github"
)

func Portal(c *gin.Context) {
	client := c.MustGet("gh_client").(*github.Client)

	user, _, err := client.Users.Get(c, "")
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.HTML(http.StatusOK, "portal.html", gin.H{
		"username": user.Login,
	})
}
