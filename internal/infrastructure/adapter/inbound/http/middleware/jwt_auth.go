package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/apex20/backend/internal/application/port"
)

type contextKey string

const claimsContextKey contextKey = "auth_claims"

// WithAuthClaims returns a new context carrying the given claims.
// Exported for use in tests that bypass the middleware.
func WithAuthClaims(ctx context.Context, claims port.AuthClaims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// ClaimsFromContext retrieves AuthClaims from the context.
func ClaimsFromContext(ctx context.Context) (port.AuthClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(port.AuthClaims)
	return claims, ok
}

// UserIDFromContext is a convenience wrapper used by handlers.
func UserIDFromContext(ctx context.Context) (port.AuthClaims, bool) {
	return ClaimsFromContext(ctx)
}

// JWTAuth returns a Chi-compatible middleware that validates Bearer tokens.
// Paths whose prefix matches any of publicPrefixes are skipped.
func JWTAuth(validator port.TokenValidator, publicPrefixes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, prefix := range publicPrefixes {
				if strings.HasPrefix(r.URL.Path, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeUnauthorized(w, "authentication required")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := validator.Validate(tokenStr)
			if err != nil {
				writeUnauthorized(w, "invalid or expired token")
				return
			}

			next.ServeHTTP(w, r.WithContext(WithAuthClaims(r.Context(), claims)))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"title":  "Unauthorized",
		"status": http.StatusUnauthorized,
		"detail": detail,
	})
}
