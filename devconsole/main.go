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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/accrescent/devconsole/config"
	"github.com/accrescent/devconsole/data"
)

func main() {
	db := new(data.SQLite)
	if err := db.Open("devconsole.db?_fk=yes&_journal=WAL"); err != nil {
		log.Fatal(err)
	}
	if err := db.Initialize(); err != nil {
		log.Fatal(err)
	}

	fileStorage := data.NewLocalStorage("/")

	oauth2Conf, conf, err := db.LoadConfig()
	if err != nil {
		oauth2Conf = &oauth2.Config{
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
		conf = &config.Config{
			SignerGitHubID: signerGitHubID,
			RepoURL:        os.Getenv("REPO_URL"),
			APIKey:         os.Getenv("API_KEY"),
		}
		log.Println("Loading config from DB failed. Falling back to environment.")
	}

	app, err := NewApp(db, fileStorage, *oauth2Conf, *conf)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := app.Start(); err != nil && errors.Is(http.ErrServerClosed, err) {
			log.Println(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Fatal("Shutting down forcefully:", err)
	}
}
