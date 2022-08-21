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

	"github.com/gin-gonic/gin"

	"github.com/accrescent/reposerver/api"
)

func main() {
	router := gin.New()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("API_KEY")

	auth := router.Group("/", authRequired(apiKey))
	auth.POST("/apps/:id/:versionCode/:version", api.PublishApp)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
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
