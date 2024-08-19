package sessionservice

import "errors"

var (
	ErrInternal            = errors.New("internal error")
	ErrNoSession           = errors.New("no session")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
