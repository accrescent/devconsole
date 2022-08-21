package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/accrescent/devportal/api"
	"github.com/accrescent/devportal/auth"
	"github.com/accrescent/devportal/config"
	"github.com/accrescent/devportal/middleware"
	"github.com/accrescent/devportal/page"
)

func main() {
	r := gin.New()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "devportal.db?_fk=yes")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		gh_id INT NOT NULL,
		access_token TEXT NOT NULL,
		expiry_time INT NOT NULL
	) STRICT`); err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		gh_id INT PRIMARY KEY,
		email TEXT NOT NULL,
		reviewer INT NOT NULL CHECK (reviewer IN (FALSE, TRUE)) DEFAULT FALSE
	) STRICT`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS usable_email_cache (
		session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
		email TEXT NOT NULL,
		PRIMARY KEY (session_id, email)
	) STRICT`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS staging_apps (
		id TEXT NOT NULL,
		session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
		path TEXT NOT NULL,
		PRIMARY KEY (id, session_id)
	) STRICT`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS review_errors (
		id TEXT PRIMARY KEY
	) STRICT`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS staging_app_review_errors (
		staging_app_id TEXT NOT NULL,
		staging_app_session_id TEXT NOT NULL,
		review_error_id TEXT NOT NULL REFERENCES review_errors(id) ON DELETE CASCADE,
		PRIMARY KEY (staging_app_id, staging_app_session_id, review_error_id),
		FOREIGN KEY (staging_app_id, staging_app_session_id)
			REFERENCES staging_apps(id, session_id)
			ON DELETE CASCADE
	) STRICT`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS submitted_apps (
		id TEXT PRIMARY KEY,
		gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
		path TEXT NOT NULL
	) STRICT`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS submitted_app_review_errors (
		submitted_app_id TEXT NOT NULL REFERENCES submitted_apps(id) ON DELETE CASCADE,
		review_error_id TEXT NOT NULL REFERENCES review_errors(id) ON DELETE CASCADE,
		PRIMARY KEY (submitted_app_id, review_error_id)
	) STRICT`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS app_teams (
		id TEXT PRIMARY KEY
	) STRICT`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS app_team_users (
		app_id TEXT NOT NULL REFERENCES app_teams(id) ON DELETE CASCADE,
		user_gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
		PRIMARY KEY (app_id, user_gh_id)
	) STRICT`); err != nil {
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

	r.LoadHTMLGlob("page/templates/*.html")

	r.Use(middleware.DB(db))
	r.Use(middleware.OAuth2Config(oauth2Conf))
	r.Use(middleware.Config(conf))

	r.GET("/auth/github", auth.GitHub)
	r.GET("/auth/github/callback", auth.GitHubCallback)
	r.StaticFile("/auth/redirect/register", "./page/static/redirect_register.html")
	r.StaticFile("/auth/redirect/dashboard", "./page/static/redirect_dashboard.html")

	auth := r.Group("/", middleware.AuthRequired())
	auth.GET("/register", page.Register)
	auth.GET("/dashboard", page.Dashboard)
	auth.StaticFile("/apps/new", "./page/static/new_app.html")
	auth.POST("/api/register", api.Register)
	auth.POST("/api/logout", api.Logout)
	auth.POST("/api/apps", api.NewApp)
	auth.PATCH("/api/apps", api.SubmitApp)
	auth.POST("/api/apps/approve", api.ApproveApp)
	auth.POST("/api/apps/:appID", middleware.SignerRequired(), api.PublishApp)

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
