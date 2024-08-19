package househandlerconverter

import (
	househandlermodel "avito/internal/handler/house/model"
	"avito/internal/model"
)

func ToHouseHandlerModel(house model.House) househandlermodel.House {
	return househandlermodel.House{
		HouseId:              house.HouseId,
		Address:              house.Address,
		Year:                 house.Year,
		Developer:            house.Developer,
		CreatedAt:            house.CreatedAt,
		LastApartmentAddedAt: house.LastApartmentAddedAt,
	}
}
