package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/accrescent/devportal/api"
	"github.com/accrescent/devportal/auth"
	"github.com/accrescent/devportal/config"
	"github.com/accrescent/devportal/data"
	"github.com/accrescent/devportal/middleware"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}

	db, err := data.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	if err := data.InitializeDB(db); err != nil {
		log.Fatal(err)
	}

	oauth2Conf := &oauth2.Config{
		ClientID:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
		Endpoint:     endpoints.GitHub,
		RedirectURL:  os.Getenv("OAUTH2_REDIRECT_URL"),
		Scopes:       []string{"user:email"},
	}

	signerGitHubID, err := strconv.ParseInt(os.Getenv("SIGNER_GH_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	conf := &config.Config{
		SignerGitHubID: signerGitHubID,
		RepoURL:        os.Getenv("REPO_URL"),
		APIKey:         os.Getenv("API_KEY"),
	}

	r.Use(middleware.DB(db))
	r.Use(middleware.OAuth2Config(oauth2Conf))
	r.Use(middleware.Config(conf))

	r.GET("/auth/github", auth.GitHub)
	r.GET("/auth/github/callback", auth.GitHubCallback)

	auth := r.Group("/", middleware.AuthRequired())
	update := auth.Group("/", middleware.UserCanUpdateRequired())
	auth.POST("/api/register", api.Register)
	auth.POST("/api/logout", api.LogOut)
	auth.POST("/api/apps", api.NewApp)
	auth.PATCH("/api/apps", api.SubmitApp)
	update.POST("/api/apps/:id/updates", api.UpdateApp)
	update.PATCH("/api/apps/:id/updates", api.SubmitAppUpdate)
	auth.POST("/api/apps/approve", middleware.ReviewerRequired(), api.ApproveApp)
	auth.POST("/api/apps/:id/updates/:version/approve", middleware.ReviewerRequired(), api.ApproveUpdate)
	auth.POST("/api/apps/:id", middleware.SignerRequired(), api.PublishApp)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(http.ErrServerClosed, err) {
			log.Println(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutting down forcefully:", err)
	}
}
