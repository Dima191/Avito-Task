package middleware

import (
	tokenmanager "avito/pkg/token_manager"
	tokenmanagerimpl "avito/pkg/token_manager/implementation"
	"errors"
	"net/http"
	"reflect"
	"slices"
)

func CheckRole(tm tokenmanager.Manager, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := parseAuthHeader(tm, r.Header.Get(AuthorizationHeader))
			if err != nil {
				switch {
				case errors.Is(err, ErrInvalidAuthHeader):
					http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
					return
				default:
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
			}

			if slices.Contains(allowedRoles, reflect.ValueOf(claims[tokenmanagerimpl.RoleClaimsTag]).String()) {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusForbidden)
		})
	}
}
