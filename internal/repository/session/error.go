package sessionrepository

import "errors"

var (
	ErrInternal  = errors.New("internal error")
	ErrNoSession = errors.New("no session")
)
