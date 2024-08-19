package apartmentserviceimpl

import (
	"avito/internal/model"
	apartmentrepository "avito/internal/repository/apartment"
	apartmentservice "avito/internal/service/apartment"
	"context"
	"errors"
	"log/slog"
)

const (
	moderator = "moderator"
)

type service struct {
	rep    apartmentrepository.Repository
	logger *slog.Logger
}

func (s *service) Create(ctx context.Context, apartment model.Apartment) error {
	if err := s.rep.Create(ctx, apartment); err != nil {
		switch {
		case errors.Is(err, apartmentrepository.ErrInvalidHouseID):
			return apartmentservice.ErrInvalidHouseID
		default:
			return apartmentservice.ErrInternal
		}
	}

	return nil
}

func (s *service) Update(ctx context.Context, apartment model.Apartment) error {
	if err := s.rep.Update(ctx, apartment); err != nil {
		switch {
		case errors.Is(err, apartmentrepository.ErrInvalidHouseID):
			return apartmentservice.ErrInvalidHouseID
		default:
			return apartmentservice.ErrInternal
		}
	}

	return nil
}

func (s *service) Apartments(ctx context.Context, houseID uint32, offset int, limit int, role string) ([]model.Apartment, error) {
	apartments, err := s.rep.Apartments(ctx, houseID, offset, limit, role != moderator)
	if err != nil {
		return nil, apartmentservice.ErrInternal
	}

	return apartments, nil
}

func New(rep apartmentrepository.Repository, logger *slog.Logger) apartmentservice.Service {
	s := &service{
		rep:    rep,
		logger: logger,
	}

	return s
}
