package server

import (
	"encoding/json"
	"net/http"

	"github.com/accrescent/devportal/dbutil"
)

func (s *Server) HandleRegisterNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var email struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	sid := r.Context().Value("sid")

	var valid bool
	if err := s.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM valid_email_cache
		WHERE session_id = ? AND email = ?
	)`, sid, email.Email).Scan(&valid); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if valid {
		if _, err := s.DB.Exec(
			"DELETE FROM valid_email_cache WHERE session_id = ?",
			sid, email.Email); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ghId, err := dbutil.GetUserID(s.DB, sid.(string))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	res, err := s.DB.Exec(
		"INSERT INTO users (gh_id, email) VALUES (?, ?)",
		ghId, email.Email,
	)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if rows != 1 {
		http.Error(w, "", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}
