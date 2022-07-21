package auth

import "errors"

var (
	ErrNoStateParam = errors.New("no state param passed in OAuth2 callback")
	ErrNoStateMatch = errors.New("state param doesn't match expected value")
)
