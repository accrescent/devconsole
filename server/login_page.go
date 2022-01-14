package server

import "net/http"

func (s *Server) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "web/static/login.html")
}
