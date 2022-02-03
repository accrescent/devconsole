package server

import (
	"net/http"

	"github.com/accrescent/devportal/dbutil"
)

func (s *Server) RegisteredMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid := r.Context().Value("sid")

		registered, err := dbutil.IsUserRegistered(s.DB, sid.(string))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if registered {
			h(w, r)
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:   "__Host-session",
				Path:   "/",
				MaxAge: -1,
				Secure: true,
			})
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}
	}
}
