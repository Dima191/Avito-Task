package apartmenthandlerconverter

import (
	apartmenthandlermodel "avito/internal/handler/apartment/model"
	"avito/internal/model"
)

func ToHandlerModelApartment(apartment model.Apartment) apartmenthandlermodel.Apartment {
	return apartmenthandlermodel.Apartment{
		ID:               apartment.ID,
		ApartmentNumber:  apartment.ApartmentNumber,
		HouseID:          apartment.HouseID,
		Price:            apartment.Price,
		NumberOfRooms:    apartment.NumberOfRooms,
		ModerationStatus: apartment.ModerationStatus,
	}
}
