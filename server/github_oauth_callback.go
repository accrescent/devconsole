package server

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v42/github"

	"github.com/accrescent/devportal/dbutil"
)

func (s *Server) HandleGitHubOAuthCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	state := r.FormValue("state")
	stateCookie, err := r.Cookie("__Host-oauth_state")
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	if stateCookie.Value != state {
		http.SetCookie(w, &http.Cookie{
			Name:   "__Host-oauth_state",
			Path:   "/",
			MaxAge: -1,
			Secure: true,
		})
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	code := r.FormValue("code")
	token, err := s.OAuth2Conf.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	sid := make([]byte, 16)
	_, err = rand.Read(sid)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	sidStr := hex.EncodeToString(sid)

	now := time.Now()
	_, err = s.DB.Exec("DELETE FROM sessions WHERE expiry_time < ?", now.Unix())
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	_, err = s.DB.Exec(
		"INSERT INTO sessions (id, gh_id, access_token, expiry_time) VALUES (?, ?, ?, ?)",
		sidStr, user.ID, token.AccessToken, now.Add(24*time.Hour).Unix(),
	)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "__Host-session",
		Value:    sidStr,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	registered, err := dbutil.IsUserRegistered(s.DB, sidStr)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if registered {
		http.Redirect(w, r, "/portal", http.StatusFound)
	} else {
		http.Redirect(w, r, "/register", http.StatusFound)
	}
}
