package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/accrescent/devconsole/api"
	"github.com/accrescent/devconsole/auth"
	"github.com/accrescent/devconsole/config"
	"github.com/accrescent/devconsole/data"
	"github.com/accrescent/devconsole/middleware"
)

type App struct {
	server      *http.Server
	db          data.DB
	fileStorage data.FileStorage
}

func NewApp(
	db data.DB,
	dsn string,
	fileStorage data.FileStorage,
	oauth2Conf oauth2.Config,
	conf config.Config,
) (*App, error) {
	if err := db.Open(dsn); err != nil {
		return nil, err
	}
	if err := db.Initialize(); err != nil {
		return nil, err
	}

	router := gin.New()
	router.Use(gin.Logger())
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}
	router.Use(middleware.DB(db))
	router.Use(middleware.FileStorage(data.NewLocalStorage("/")))
	router.Use(middleware.OAuth2Config(oauth2Conf))
	router.Use(middleware.Config(conf))

	router.GET("/auth/github", auth.GitHub)
	router.GET("/api/auth/github/callback", auth.GitHubCallback)

	auth := router.Group("/", middleware.AuthRequired())
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

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return &App{
		server,
		db,
		fileStorage,
	}, nil
}

func (a *App) Start() error {
	return a.server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	return a.db.Close()
}
