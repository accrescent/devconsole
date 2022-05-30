package server

import (
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v45/github"
)

var regTmpl = template.Must(template.New("register.html").ParseFiles("web/templates/register.html"))

// The format for GitHub's noreply email addresses is documented at
// https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-github-user-account/managing-email-preferences/setting-your-commit-email-address
var noReplyEmail = regexp.MustCompile(`^([0-9]{7}\+)?.*@users\.noreply\.github\.com$`)

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	sid := ctx.Value("sid")

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
		http.Redirect(w, r, "/portal", http.StatusFound)
	}

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

	// ListOptions are the defaults from
	// https://docs.github.com/en/rest/reference/users#list-email-addresses-for-the-authenticated-user
	emails, _, err := client.Users.ListEmails(ctx, &github.ListOptions{Page: 1, PerPage: 30})
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	usableEmails := []string{}
	for _, email := range emails {
		address := email.GetEmail()
		if email.GetVerified() && !noReplyEmail.MatchString(address) {
			usableEmails = append(usableEmails, address)
		}
	}

	emailInsertQuery := "INSERT INTO valid_email_cache (session_id, email) VALUES "
	inserts := []string{}
	params := []interface{}{}
	for _, email := range usableEmails {
		inserts = append(inserts, "(?, ?)")
		params = append(params, sid, email)
	}
	emailInsertQuery = emailInsertQuery + strings.Join(inserts, ",")
	if _, err = s.DB.Exec(emailInsertQuery, params...); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	registerData := struct {
		Username *string
		Emails   []string
	}{user.Login, usableEmails}

	regTmpl.Execute(w, registerData)
}
