package apartmentrepository

import (
	"avito/internal/model"
	"context"
)

type Repository interface {
	Create(ctx context.Context, apartment model.Apartment) error
	Update(ctx context.Context, apartment model.Apartment) error
	Apartments(ctx context.Context, houseID uint32, offset int, limit int, moderationStatusConstraint bool) ([]model.Apartment, error)
	CloseConnection() error
}
