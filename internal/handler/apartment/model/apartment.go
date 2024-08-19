package apartmenthandlermodel

import (
	"github.com/go-playground/validator/v10"
	"slices"
)

type Apartment struct {
	ID               uint32 `json:"id,omitempty"`
	ApartmentNumber  int    `json:"apartmentNumber"`
	HouseID          uint32 `json:"house_id" validate:"required"`
	Price            uint32 `json:"price" validate:"required"`
	NumberOfRooms    uint32 `json:"number_of_rooms" validate:"required"`
	ModerationStatus string `json:"moderation_status" validate:"moderation_status"`
}

var (
	PossibleModerationStatus = []string{"created", "approved", "declined", "on moderation"}
)

func ModerationStatusValidation(fl validator.FieldLevel) bool {
	return slices.Contains(PossibleModerationStatus, fl.Field().String())
}
