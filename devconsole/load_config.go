//go:build !debug

package main

import (
	"golang.org/x/oauth2"

	"github.com/accrescent/devconsole/config"
	"github.com/accrescent/devconsole/data"
)

func loadConfig(db data.DB) (*oauth2.Config, *config.Config, error) {
	return db.LoadConfig()
}
