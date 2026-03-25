package middleware_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	jwtinfra "github.com/apex20/backend/internal/infrastructure/adapter/outbound/jwt"
	"github.com/apex20/backend/internal/infrastructure/adapter/inbound/http/middleware"
)

func setup(t *testing.T) (*rsa.PrivateKey, port.TokenValidator) {
	t.Helper()
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return pk, jwtinfra.NewRSATokenValidator(&pk.PublicKey)
}

func okHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.ClaimsFromContext(r.Context())
		assert.True(t, ok, "expected claims in context")
		assert.NotEqual(t, uuid.Nil, claims.UserID)
		w.WriteHeader(http.StatusOK)
	})
}

func TestJWTAuth_ValidToken_PassesAndSetsContext(t *testing.T) {
	pk, validator := setup(t)
	userID := uuid.New()

	gen := jwtinfra.NewRSATokenGenerator(pk, time.Hour)
	token, err := gen.Generate(userID, true)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/campaigns", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	middleware.JWTAuth(validator)(okHandler(t)).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestJWTAuth_MissingHeader_Returns401(t *testing.T) {
	_, validator := setup(t)

	req := httptest.NewRequest(http.MethodGet, "/campaigns", nil)
	rr := httptest.NewRecorder()

	middleware.JWTAuth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTAuth_InvalidToken_Returns401(t *testing.T) {
	_, validator := setup(t)

	req := httptest.NewRequest(http.MethodGet, "/campaigns", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rr := httptest.NewRecorder()

	middleware.JWTAuth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTAuth_NonBearerScheme_Returns401(t *testing.T) {
	_, validator := setup(t)

	req := httptest.NewRequest(http.MethodGet, "/campaigns", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	rr := httptest.NewRecorder()

	middleware.JWTAuth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTAuth_PublicPath_SkipsValidation(t *testing.T) {
	_, validator := setup(t)

	publicPaths := []string{"/health", "/apex20.v1.AuthService/"}

	for _, path := range publicPaths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()

			called := false
			middleware.JWTAuth(validator, "/health", "/apex20.v1.AuthService/")(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					called = true
					w.WriteHeader(http.StatusOK)
				}),
			).ServeHTTP(rr, req)

			assert.True(t, called)
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}

func TestJWTAuth_ExpiredToken_Returns401(t *testing.T) {
	pk, validator := setup(t)

	gen := jwtinfra.NewRSATokenGenerator(pk, -time.Hour)
	token, err := gen.Generate(uuid.New(), false)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/campaigns", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rr := httptest.NewRecorder()

	middleware.JWTAuth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestWithAuthClaims_SetsAndReadsContext(t *testing.T) {
	userID := uuid.New()
	claims := port.AuthClaims{UserID: userID, IsAdmin: true}

	ctx := middleware.WithAuthClaims(context.Background(), claims)
	got, ok := middleware.ClaimsFromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, claims, got)
}
