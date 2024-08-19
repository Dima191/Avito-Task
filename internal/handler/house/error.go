package househandler

import "errors"

var (
	ErrHouseAlreadyExists = errors.New("house with this address already exists")
	ErrHouseNotFound      = errors.New("house not found")
)
