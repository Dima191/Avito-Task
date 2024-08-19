package apartmentrepositoryconverter

import (
	"avito/internal/model"
	apartmentrepositorymodel "avito/internal/repository/apartment/model"
)

func ToApartmentDTO(apartment apartmentrepositorymodel.Apartment) model.Apartment {
	return model.Apartment{
		ID:               apartment.ID,
		ApartmentNumber:  apartment.ApartmentNumber,
		HouseID:          apartment.HouseID,
		Price:            apartment.Price,
		NumberOfRooms:    apartment.NumberOfRooms,
		ModerationStatus: apartment.ModerationStatus,
	}
}
