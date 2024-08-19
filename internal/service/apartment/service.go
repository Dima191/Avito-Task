package apartmentservice

import (
	"avito/internal/model"
	"context"
)

type Service interface {
	Create(ctx context.Context, apartment model.Apartment) error
	Update(ctx context.Context, apartment model.Apartment) error
	Apartments(ctx context.Context, houseID uint32, offset int, limit int, role string) ([]model.Apartment, error)
}
