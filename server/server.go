package server

import (
	"net/http"

	"golang.org/x/oauth2"
)

type Server struct {
	Router     *http.ServeMux
	OAuth2Conf oauth2.Config
	Sessions   map[string]string
}
