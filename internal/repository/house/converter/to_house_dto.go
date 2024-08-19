package houserepositoryconverter

import (
	"avito/internal/model"
	houserepositorymodel "avito/internal/repository/house/model"
)

func ToHouseDto(house houserepositorymodel.House) model.House {
	return model.House{
		HouseId:              house.HouseId,
		Address:              house.Address,
		Year:                 house.Year,
		Developer:            house.Developer,
		CreatedAt:            house.CreatedAt,
		LastApartmentAddedAt: house.LastApartmentAddedAt,
	}
}
