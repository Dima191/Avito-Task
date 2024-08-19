package userrepository

import "errors"

var (
	ErrInternal          = errors.New("internal error")
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrUserNotFound      = errors.New("user not found")
)
