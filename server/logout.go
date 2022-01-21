package server

import "net/http"

func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	sid, err := r.Cookie("__Host-session")
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	_, err = s.DB.Exec("DELETE FROM sessions WHERE id = ?", sid.Value)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "__Host-session",
		Path:   "/",
		MaxAge: -1,
		Secure: true,
	})

	w.WriteHeader(http.StatusOK)
}
