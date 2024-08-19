package houseservice

import (
	"avito/internal/model"
	"context"
)

type Service interface {
	Create(ctx context.Context, house model.House) error
	Houses(ctx context.Context, offset int, limit int) ([]model.House, error)
}
