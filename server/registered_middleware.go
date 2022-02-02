package server

import "net/http"

func (s *Server) RegisteredMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid := r.Context().Value("sid")
		var ghId string
		if err := s.DB.QueryRow(
			"SELECT gh_id FROM sessions WHERE id = ?",
			sid,
		).Scan(&ghId); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		var registered bool
		if err := s.DB.QueryRow(
			"SELECT EXISTS (SELECT 1 FROM users WHERE gh_id = ?)",
			ghId,
		).Scan(&registered); err != nil {
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
