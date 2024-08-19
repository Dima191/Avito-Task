package userhandler

import "errors"

var (
	ErrDecodeBody          = errors.New("failed to decode request body")
	ErrCredentialsInvalid  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidURLParams    = errors.New("invalid url params")
	ErrEmailAlreadyTaken   = errors.New("email already taken")
	ErrNoSession           = errors.New("no session by refresh token. login again")
)
