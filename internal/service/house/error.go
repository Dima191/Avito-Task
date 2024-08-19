package houseservice

import "errors"

var (
	ErrInternal           = errors.New("internal server error")
	ErrHouseAlreadyExists = errors.New("house with this address already exists")
	ErrHouseNotFound      = errors.New("house not found")
)
