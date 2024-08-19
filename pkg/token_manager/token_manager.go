package tokenmanager

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Manager interface {
	GenerateAccessToken(userID uint32, role string) (accessToken string, err error)
	GenerateRefreshToken() (refreshToken string, expiresAt time.Time, err error)
	Parse(tokenString string) (jwt.MapClaims, error)
}
