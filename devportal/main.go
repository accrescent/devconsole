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

	db := new(data.SQLite)
	if err := db.Open("devportal.db?_fk=yes&_journal=WAL"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Initialize(); err != nil {
		log.Fatal(err)
	}

	oauth2Conf := oauth2.Config{
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
	conf := config.Config{
		SignerGitHubID: signerGitHubID,
		RepoURL:        os.Getenv("REPO_URL"),
		APIKey:         os.Getenv("API_KEY"),
	}

	r.Use(middleware.DB(db))
	r.Use(middleware.FileStorage(data.NewLocalStorage("/")))
	r.Use(middleware.OAuth2Config(oauth2Conf))
	r.Use(middleware.Config(conf))

	r.GET("/auth/github", auth.GitHub)
	r.GET("/api/auth/github/callback", auth.GitHubCallback)

	auth := r.Group("/", middleware.AuthRequired())
	reviewer := auth.Group("/", middleware.ReviewerRequired())
	update := auth.Group("/", middleware.UserCanUpdateRequired())
	auth.GET("/api/emails", api.GetEmails)
	reviewer.GET("/api/pending-apps", api.GetPendingApps)
	reviewer.GET("/api/pending-apps/:id/apks", api.GetAppAPKs)
	reviewer.PATCH("/api/pending-apps/:id", api.ApproveApp)
	reviewer.DELETE("/api/pending-apps/:id", api.RejectApp)
	reviewer.GET("/api/updates", api.GetUpdates)
	reviewer.GET("/api/updates/:id/:version/apks", api.GetUpdateAPKs)
	reviewer.PATCH("/api/updates/:id/:version", api.ApproveUpdate)
	reviewer.DELETE("/api/updates/:id/:version", api.RejectUpdate)
	auth.GET("/api/approved-apps", middleware.SignerRequired(), api.GetApprovedApps)
	auth.POST("/api/register", api.Register)
	auth.DELETE("/api/session", api.LogOut)
	auth.GET("/api/apps", api.GetApps)
	auth.POST("/api/apps", api.NewApp)
	auth.PATCH("/api/apps/:id", api.SubmitApp)
	update.POST("/api/apps/:id/updates", api.NewUpdate)
	update.PATCH("/api/apps/:id/:version", api.SubmitUpdate)
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
