package server

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

func (s *Server) HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	state := make([]byte, 16)
	_, err := rand.Read(state)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	stateStr := hex.EncodeToString(state)

	http.SetCookie(w, &http.Cookie{
		Name:     "__Host-oauth_state",
		Value:    stateStr,
		Path:     "/",
		MaxAge:   30,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	url := s.OAuth2Conf.AuthCodeURL(stateStr)
	http.Redirect(w, r, url, http.StatusFound)
}
