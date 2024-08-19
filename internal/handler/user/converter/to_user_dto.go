package userhandlerconverter

import (
	userhandlermodel "avito/internal/handler/user/model"
	"avito/internal/model"
	"avito/pkg/hasher"
	"github.com/google/uuid"
)

func ToUserDto(user userhandlermodel.User) (model.User, error) {
	hash, err := hasher.Hash(user.Password)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:           uuid.New().ID(),
		Role:         user.Role,
		Email:        user.Email,
		HashPassword: hash,
	}, nil
}
