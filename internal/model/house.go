package model

import "time"

type House struct {
	HouseId              uint32
	Address              string
	Year                 int
	Developer            string
	CreatedAt            time.Time
	LastApartmentAddedAt time.Time
}
