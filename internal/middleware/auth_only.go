package middleware

import (
	tokenmanager "avito/pkg/token_manager"
	tokenmanagerimpl "avito/pkg/token_manager/implementation"
	"context"
	"errors"
	"net/http"
	"reflect"
	"time"
)

func AuthOnly(tm tokenmanager.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := parseAuthHeader(tm, r.Header.Get(AuthorizationHeader))
			if err != nil {
				switch {
				case errors.Is(err, ErrInvalidAuthHeader):
					http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
					return
				case errors.Is(err, ErrInvalidToken):
					http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
					return
				default:
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
			}

			expAt, err := claims.GetExpirationTime()
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if time.Now().After(expAt.Time) {
				w.WriteHeader(http.StatusUnauthorized)
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
