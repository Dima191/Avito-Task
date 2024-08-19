package tokenmanagerimpl

import (
	tokenmanager "avito/pkg/token_manager"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"unsafe"
)

const (
	RoleClaimsTag   = "role"
	ExpClaimsTag    = "exp"
	UserIDClaimsTag = "user_id"
)

type manager struct {
	jwtSignedKey          []byte
	accessTokenExpiresIn  time.Duration
	refreshTokenExpiresIn time.Duration
}

func (m *manager) GenerateAccessToken(userID uint32, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		UserIDClaimsTag: userID,
		RoleClaimsTag:   role,
		ExpClaimsTag:    time.Now().Add(m.accessTokenExpiresIn).Unix(),
	})

	return token.SignedString(m.jwtSignedKey)
}

func (m *manager) GenerateRefreshToken() (refreshToken string, expiresAt time.Time, err error) {
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		return "", time.Time{}, err
	}

	return hex.EncodeToString(token), time.Now().Add(m.refreshTokenExpiresIn), nil
}

func (m *manager) Parse(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.jwtSignedKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
		default:
			return nil, err
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func New(jwtSignedKey string, accessTokenExpiresIn, refreshTokenExpiresIn time.Duration) tokenmanager.Manager {
	tm := &manager{
		jwtSignedKey:          unsafe.Slice(unsafe.StringData(jwtSignedKey), len(jwtSignedKey)),
		accessTokenExpiresIn:  accessTokenExpiresIn,
		refreshTokenExpiresIn: refreshTokenExpiresIn,
	}

	return tm
}
