package server

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/accrescent/devportal/dbutil"
)

var portalTmpl = template.Must(template.New("portal.html").ParseFiles("web/templates/portal.html"))

func (s *Server) HandlePortal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	sid := ctx.Value("sid")

	client, err := dbutil.CreateGitHubClient(s.DB, sid.(string), ctx)
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

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	portalTmpl.Execute(w, user.Login)
}
