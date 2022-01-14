package server

import (
	"html/template"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v41/github"
)

var tmpl = template.Must(template.New("portal.html").ParseFiles("web/templates/portal.html"))

func (s *Server) HandlePortal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	token := r.Context().Value("token")
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.(string)})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	tmpl.Execute(w, user.Login)
}
