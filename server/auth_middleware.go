package server

import (
	"context"
	"net/http"
	"time"
)

func (s *Server) AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, err := r.Cookie("__Host-session")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		var exists bool
		if err = s.DB.QueryRow(
			"SELECT EXISTS (SELECT 1 FROM sessions WHERE id = ? AND expiry_time > ?)",
			sid.Value, time.Now().Unix(),
		).Scan(&exists); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.SetCookie(w, &http.Cookie{
				Name:   "__Host-session",
				Path:   "/",
				MaxAge: -1,
				Secure: true,
			})
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "sid", sid.Value)

		h(w, r.WithContext(ctx))
	}
}
