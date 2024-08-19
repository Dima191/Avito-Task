package userrepositoryconverter

import (
	"avito/internal/model"
	userrepositorymodel "avito/internal/repository/user/model"
)

func ToUserDTO(user userrepositorymodel.User) model.User {
	return model.User{
		ID:           user.ID,
		Role:         user.Role,
		Email:        user.Email,
		HashPassword: user.HashPassword,
	}
}
