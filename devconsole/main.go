package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/accrescent/devconsole/data"
)

//go:generate protoc -I proto --go_out pb proto/commands.proto proto/config.proto proto/targeting.proto

func main() {
	setMode()

	db := new(data.SQLite)
	if err := db.Open("devconsole.db?_fk=yes&_journal=WAL"); err != nil {
		log.Fatal(err)
	}
	if err := db.Initialize(); err != nil {
		log.Fatal(err)
	}

	fileStorage := data.NewLocalStorage(".")

	oauth2Conf, conf, err := loadConfig(db)
	if err != nil {
		log.Fatal(err)
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
