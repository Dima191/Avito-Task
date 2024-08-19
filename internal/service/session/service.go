package sessionservice

import (
	"context"
)

type Service interface {
	Create(ctx context.Context, userID uint32, role string) (accessToken, refreshToken string, err error)
	Update(ctx context.Context, userID uint32, role string, expiredRefreshToken string) (accessToken, refreshToken string, err error)
	ResetSession(ctx context.Context, userID uint32, role string) (accessToken, refreshToken string, err error)
}
