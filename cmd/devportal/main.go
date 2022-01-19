package main

import (
	"log"
	"net/http"
	"os"
	"time"

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
	s := server.Server{
		Router:     http.NewServeMux(),
		OAuth2Conf: conf,
		Sessions:   make(map[string]string),
	}

	s.Router.HandleFunc("/auth/github", s.HandleGitHubLogin)
	s.Router.HandleFunc("/auth/github/callback", s.HandleGitHubOAuthCallback)
	s.Router.HandleFunc("/portal", s.AuthMiddleware(s.HandlePortal))

	http := &http.Server{
		Addr:         ":8080",
		Handler:      s.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(http.ListenAndServe())
}
