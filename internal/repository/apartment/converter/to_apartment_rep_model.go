package apartmentrepositoryconverter

import (
	"avito/internal/model"
	apartmentrepositorymodel "avito/internal/repository/apartment/model"
)

func ToApartmentRepModel(apartment model.Apartment) apartmentrepositorymodel.Apartment {
	return apartmentrepositorymodel.Apartment{
		ID:               apartment.ID,
		ApartmentNumber:  apartment.ApartmentNumber,
		HouseID:          apartment.HouseID,
		Price:            apartment.Price,
		NumberOfRooms:    apartment.NumberOfRooms,
		ModerationStatus: apartment.ModerationStatus,
	}
}
