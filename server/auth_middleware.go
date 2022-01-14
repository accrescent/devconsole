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

		token, exists := s.Sessions[sid.Value]
		if !exists {
			http.SetCookie(w, &http.Cookie{
				Name:   "__Host-session",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)

		h(w, r.WithContext(ctx))
	}
}
