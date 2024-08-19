package houserepositoryconverter

import (
	"avito/internal/model"
	houserepositorymodel "avito/internal/repository/house/model"
)

func ToHouseRepModel(house model.House) houserepositorymodel.House {
	return houserepositorymodel.House{
		HouseId:              house.HouseId,
		Address:              house.Address,
		Year:                 house.Year,
		Developer:            house.Developer,
		CreatedAt:            house.CreatedAt,
		LastApartmentAddedAt: house.LastApartmentAddedAt,
	}
}
