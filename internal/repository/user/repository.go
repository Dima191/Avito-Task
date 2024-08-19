package userrepository

import (
	"avito/internal/model"
	"context"
)

type Repository interface {
	Save(ctx context.Context, user model.User) error
	UserByEmail(ctx context.Context, email string) (model.User, error)
	CloseConnection() error
}
