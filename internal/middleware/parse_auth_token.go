package middleware

import (
	tokenmanager "avito/pkg/token_manager"
	tokenmanagerimpl "avito/pkg/token_manager/implementation"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"reflect"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	BearerToken         = "Bearer"
)

const (
	RoleCtxKey   = "role"
	UserIDCtxKey = "user_id"
)

var (
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrInvalidToken      = errors.New("invalid token. login again or update tokens")
)

func parseAuthHeader(tm tokenmanager.Manager, authorizationToken string) (jwt.MapClaims, error) {
	authorizationTokenParts := strings.Split(authorizationToken, " ")
	if len(authorizationTokenParts) != 2 && authorizationTokenParts[0] != BearerToken {
		return map[string]interface{}{}, ErrInvalidAuthHeader
	}

	accessToken := authorizationTokenParts[1]

	claims, err := tm.Parse(accessToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return map[string]interface{}{}, ErrInvalidToken
		default:
			return map[string]interface{}{}, ErrInvalidAuthHeader
		}
	}

	return claims, nil
}

func ParseAuthToken(tm tokenmanager.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := parseAuthHeader(tm, r.Header.Get(AuthorizationHeader))
			if errors.Is(err, ErrInvalidAuthHeader) {
				http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
				return
			}

			userID := uint32(reflect.ValueOf(claims[tokenmanagerimpl.UserIDClaimsTag]).Float())
			ctx := context.WithValue(r.Context(), UserIDCtxKey, userID)

			ctx = context.WithValue(ctx, RoleCtxKey, claims[tokenmanagerimpl.RoleClaimsTag])

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
