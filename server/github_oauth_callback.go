package server

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
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
			MaxAge: -1,
		})
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	code := r.FormValue("code")
	token, err := s.OAuth2Conf.Exchange(r.Context(), code)
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

	_, err = s.DB.Exec(
		"INSERT INTO sessions (id, access_token) VALUES (?, ?)",
		sidStr, token.AccessToken,
	)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "__Host-session",
		Value:    sidStr,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/portal", http.StatusFound)
}
