package userservice

import "errors"

var (
	ErrCredentialsInvalid = errors.New("invalid email or password")
	ErrInternal           = errors.New("internal server error")
	ErrEmailAlreadyTaken  = errors.New("email already taken")
)
