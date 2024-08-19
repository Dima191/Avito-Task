package houseserviceimpl

import (
	"avito/internal/model"
	houserepository "avito/internal/repository/house"
	houseservice "avito/internal/service/house"
	"context"
	"errors"
	"log/slog"
)

type service struct {
	rep houserepository.Repository

	logger *slog.Logger
}

func (s *service) Create(ctx context.Context, house model.House) error {
	if err := s.rep.Create(ctx, house); err != nil {
		switch {
		case errors.Is(err, houserepository.ErrHouseAlreadyExists):
			return houseservice.ErrHouseAlreadyExists
		default:
			return houseservice.ErrInternal
		}
	}
	return nil
}

func (s *service) Houses(ctx context.Context, offset int, limit int) ([]model.House, error) {
	houses, err := s.rep.Houses(ctx, offset, limit)
	if err != nil {
		switch {
		case errors.Is(err, houserepository.ErrHouseNotFound):
			return nil, houseservice.ErrHouseNotFound
		default:
			return nil, houseservice.ErrInternal
		}
	}

	return houses, nil
}

func New(rep houserepository.Repository, logger *slog.Logger) houseservice.Service {
	s := &service{
		rep:    rep,
		logger: logger,
	}
	return s
}
