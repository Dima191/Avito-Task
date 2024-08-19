package userrepositoryconverter

import (
	"avito/internal/model"
	userrepositorymodel "avito/internal/repository/user/model"
)

func ToUserRepModel(user model.User) userrepositorymodel.User {
	return userrepositorymodel.User{
		ID:           user.ID,
		Role:         user.Role,
		Email:        user.Email,
		HashPassword: user.HashPassword,
	}
}
