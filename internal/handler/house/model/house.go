package househandlermodel

import "time"

type House struct {
	HouseId              uint32    `json:"house_id"`
	Address              string    `json:"address"`
	Year                 int       `json:"year"`
	Developer            string    `json:"developer"`
	CreatedAt            time.Time `json:"created_at"`
	LastApartmentAddedAt time.Time `json:"last_apartment_added_at"`
}
