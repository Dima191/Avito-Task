package userservice

import (
	"avito/internal/model"
	"context"
)

type Service interface {
	Save(ctx context.Context, user model.User) error
	LogIn(ctx context.Context, email, password string) (userID uint32, err error)
}
