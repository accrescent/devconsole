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
	}
	db, err := sql.Open("sqlite3", "devportal.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY NOT NULL,
		access_token TEXT NOT NULL,
		expiry_time INT NOT NULL
	)`)
	if err != nil {
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
	s.Router.HandleFunc("/portal", s.AuthMiddleware(s.HandlePortal))

	http := &http.Server{
		Addr:         ":8080",
		Handler:      s.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(http.ListenAndServe())
}
