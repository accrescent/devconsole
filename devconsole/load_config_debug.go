//go:build debug

package main

import (
	"os"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/accrescent/devconsole/config"
	"github.com/accrescent/devconsole/data"
)

func loadConfig(db data.DB) (*oauth2.Config, *config.Config, error) {
	oauth2Conf := &oauth2.Config{
		ClientID:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
		Endpoint:     endpoints.GitHub,
		RedirectURL:  os.Getenv("OAUTH2_REDIRECT_URL"),
		Scopes:       []string{"user:email"},
	}
	signerGitHubID, err := strconv.ParseInt(os.Getenv("SIGNER_GH_ID"), 10, 64)
	if err != nil {
		return nil, nil, err
	}
	conf := &config.Config{
		SignerGitHubID: signerGitHubID,
		RepoURL:        os.Getenv("REPO_URL"),
		APIKey:         os.Getenv("API_KEY"),
	}

	return oauth2Conf, conf, nil
}
