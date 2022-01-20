package server

import (
	"database/sql"
	"net/http"

	"golang.org/x/oauth2"
)

type Server struct {
	Router     *http.ServeMux
	OAuth2Conf oauth2.Config
	DB         *sql.DB
}
