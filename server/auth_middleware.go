package server

import (
	"context"
	"net/http"
)

func (s *Server) AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, err := r.Cookie("__Host-session")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		var token string
		err = s.DB.QueryRow("SELECT access_token FROM sessions WHERE id = ?", sid.Value).Scan(&token)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "__Host-session",
				Path: "/",
				MaxAge: -1,
				Secure: true,
			})
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)

		h(w, r.WithContext(ctx))
	}
}
