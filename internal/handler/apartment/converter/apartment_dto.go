package apartmenthandlerconverter

import (
	apartmenthandlermodel "avito/internal/handler/apartment/model"
	"avito/internal/model"
	"github.com/google/uuid"
)

func ToApartmentDTO(apartment apartmenthandlermodel.Apartment) model.Apartment {
	return model.Apartment{
		ID:               uuid.New().ID(),
		ApartmentNumber:  apartment.ApartmentNumber,
		HouseID:          apartment.HouseID,
		Price:            apartment.Price,
		NumberOfRooms:    apartment.NumberOfRooms,
		ModerationStatus: apartment.ModerationStatus,
	}
}
