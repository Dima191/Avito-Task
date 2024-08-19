package houserepository

import (
	"avito/internal/model"
	"context"
)

type Repository interface {
	Create(ctx context.Context, house model.House) error
	Houses(ctx context.Context, offset int, limit int) ([]model.House, error)
	CloseConnection() error
}
