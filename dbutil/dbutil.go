package dbutil

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v42/github"
)

// CreateGitHubClient returns a github.Client created with the OAuth2 token associated with the
// given session.
func CreateGitHubClient(db *sql.DB, sessionID string, ctx context.Context) (*github.Client, error) {
	var token string
	if err := db.QueryRow("SELECT access_token FROM sessions WHERE id = ? AND expiry_time > ?",
		sessionID, time.Now().Unix(),
	).Scan(&token); err != nil {
		return nil, err
	} else {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		return github.NewClient(tc), nil
	}
}

// GetUserID returns the GitHub user ID associated with a given session. This value corresponds to
// the "id" field from https://docs.github.com/en/rest/reference/users.
func GetUserID(db *sql.DB, sessionID string) (string, error) {
	var ghId string
	if err := db.QueryRow(
		"SELECT gh_id FROM sessions WHERE id = ?",
		sessionID,
	).Scan(&ghId); err != nil {
		return "", err
	} else {
		return ghId, nil
	}
}

// IsUserRegistered returns whether the user associated with the given session ID is registered. It
// will only return an error when it encounters a fatal database error.
func IsUserRegistered(db *sql.DB, sessionID string) (bool, error) {
	var registered bool
	if err := db.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM users WHERE gh_id = (SELECT gh_id FROM sessions WHERE id = ?
	)`, sessionID).Scan(&registered); err != nil {
		return false, err
	} else {
		return registered, nil
	}
}
