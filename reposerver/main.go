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

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"

	"github.com/accrescent/reposerver/api"
	"github.com/accrescent/reposerver/middleware"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}

	var conf config

	confFile, err := os.ReadFile("/etc/reposerver/config.toml")
	if err == nil {
		if err := toml.Unmarshal(confFile, &conf); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Loading config from file failed. Falling back to environment.")

		conf.APIKey = os.Getenv("API_KEY")
		conf.PublishDir = os.Getenv("PUBLISH_DIR")
	}

	auth := router.Group("/", middleware.AuthRequired(conf.APIKey))
	auth.Use(middleware.PublishDir(conf.PublishDir))
	auth.POST("/api/apps/:id/:versionCode/:version", api.PublishApp)
	auth.PUT("/api/apps/:id/:versionCode/:version", api.UpdateApp)

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
