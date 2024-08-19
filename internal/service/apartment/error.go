package apartmentservice

import "errors"

var (
	ErrInternal       = errors.New("internal server error")
	ErrInvalidHouseID = errors.New("invalid house id")
)
