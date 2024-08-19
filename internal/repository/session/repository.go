package sessionrepository

import (
	"avito/internal/model"
	"context"
)

type Repository interface {
	Create(ctx context.Context, session model.Session) error
	SessionByUserId(ctx context.Context, userID uint32) (session model.Session, err error)
	CheckSessionByUserId(ctx context.Context, userID uint32) (bool, error)
	ResetSession(ctx context.Context, session model.Session) error
	CloseConnection() error
}
