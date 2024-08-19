package houserepository

import "errors"

var (
	ErrInternal           = errors.New("internal error")
	ErrHouseAlreadyExists = errors.New("house with this address already exists")
	ErrHouseNotFound      = errors.New("house not found")
)
