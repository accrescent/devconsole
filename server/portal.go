package server

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v42/github"
)

var portalTmpl = template.Must(template.New("portal.html").ParseFiles("web/templates/portal.html"))

func (s *Server) HandlePortal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	sid := ctx.Value("sid")
	var token string
	err := s.DB.QueryRow(
		"SELECT access_token FROM sessions WHERE id = ? AND expiry_time > ?",
		sid, time.Now().Unix(),
	).Scan(&token)
	switch {
	case err == sql.ErrNoRows:
		http.SetCookie(w, &http.Cookie{
			Name:   "__Host-session",
			Path:   "/",
			MaxAge: -1,
			Secure: true,
		})
		http.Redirect(w, r, "/", http.StatusFound)
		return
	case err != nil:
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	portalTmpl.Execute(w, user.Login)
}
