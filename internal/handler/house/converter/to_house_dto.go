package househandlerconverter

import (
	househandlermodel "avito/internal/handler/house/model"
	"avito/internal/model"
	"github.com/google/uuid"
)

func ToHouseDTO(house househandlermodel.House) model.House {
	return model.House{
		HouseId:   uuid.New().ID(),
		Address:   house.Address,
		Year:      house.Year,
		Developer: house.Developer,
	}
}
