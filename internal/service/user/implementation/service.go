package userserviceimpl

import (
	"avito/internal/model"
	userrepository "avito/internal/repository/user"
	userservice "avito/internal/service/user"
	"avito/pkg/hasher"
	tokenmanager "avito/pkg/token_manager"
	"context"
	"errors"
	"log/slog"
)

type service struct {
	repository userrepository.Repository

	tokenManager tokenmanager.Manager

	logger *slog.Logger
}

func (s *service) Save(ctx context.Context, user model.User) error {
	if err := s.repository.Save(ctx, user); err != nil {
		switch {
		case errors.Is(err, userrepository.ErrEmailAlreadyTaken):
			return userservice.ErrEmailAlreadyTaken
		default:
			return userservice.ErrInternal
		}
	}

	return nil
}

func (s *service) LogIn(ctx context.Context, email string, password string) (userID uint32, err error) {
	user, err := s.repository.UserByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, userrepository.ErrUserNotFound):
			return 0, userservice.ErrCredentialsInvalid
		default:
			return 0, userservice.ErrInternal
		}
	}

	if err = hasher.Compare(password, user.HashPassword); err != nil {
		return 0, userservice.ErrCredentialsInvalid
	}

	return user.ID, nil
}

func New(repository userrepository.Repository, tokenManager tokenmanager.Manager, logger *slog.Logger) userservice.Service {
	s := &service{
		repository:   repository,
		tokenManager: tokenManager,
		logger:       logger,
	}
	return s
}
