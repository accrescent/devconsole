package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/accrescent/devportal/server"
)

func main() {
	conf := oauth2.Config{
		ClientID:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}
	db, err := sql.Open("sqlite3", "devportal.db?_fk=yes")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY NOT NULL,
		gh_id TEXT NOT NULL,
		access_token TEXT NOT NULL,
		expiry_time INT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		gh_id TEXT PRIMARY KEY NOT NULL,
		email TEXT NOT NULL
	)`); err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS valid_email_cache (
		session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
		email TEXT NOT NULL,
		PRIMARY KEY (session_id, email)
	)`); err != nil {
		log.Fatal(err)
	}

	s := server.Server{
		Router:     http.NewServeMux(),
		OAuth2Conf: conf,
		DB:         db,
	}

	s.Router.HandleFunc("/auth/github", s.HandleGitHubLogin)
	s.Router.HandleFunc("/auth/github/callback", s.HandleGitHubOAuthCallback)
	s.Router.HandleFunc("/logout", s.HandleLogout)
	s.Router.HandleFunc("/portal", s.AuthMiddleware(s.RegisteredMiddleware(s.HandlePortal)))
	s.Router.HandleFunc("/register", s.AuthMiddleware(s.HandleRegister))
	s.Router.HandleFunc("/register/new", s.AuthMiddleware(s.HandleRegisterNew))

	http := &http.Server{
		Addr:         ":8080",
		Handler:      s.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(http.ListenAndServe())
}
