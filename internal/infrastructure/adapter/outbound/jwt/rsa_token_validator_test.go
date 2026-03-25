package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwtinfra "github.com/apex20/backend/internal/infrastructure/adapter/outbound/jwt"
)

func TestRSATokenValidator_Validate_ValidToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	userID := uuid.New()
	gen := jwtinfra.NewRSATokenGenerator(privateKey, time.Hour)
	tokenStr, err := gen.Generate(userID, false)
	require.NoError(t, err)

	validator := jwtinfra.NewRSATokenValidator(&privateKey.PublicKey)
	claims, err := validator.Validate(tokenStr)

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.False(t, claims.IsAdmin)
}

func TestRSATokenValidator_Validate_AdminClaim(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	userID := uuid.New()
	gen := jwtinfra.NewRSATokenGenerator(privateKey, time.Hour)
	tokenStr, err := gen.Generate(userID, true)
	require.NoError(t, err)

	validator := jwtinfra.NewRSATokenValidator(&privateKey.PublicKey)
	claims, err := validator.Validate(tokenStr)

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.True(t, claims.IsAdmin)
}

func TestRSATokenValidator_Validate_ExpiredToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	gen := jwtinfra.NewRSATokenGenerator(privateKey, -time.Hour)
	tokenStr, err := gen.Generate(uuid.New(), false)
	require.NoError(t, err)

	validator := jwtinfra.NewRSATokenValidator(&privateKey.PublicKey)
	_, err = validator.Validate(tokenStr)

	assert.Error(t, err)
}

func TestRSATokenValidator_Validate_WrongPublicKey(t *testing.T) {
	privateKey1, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	privateKey2, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	gen := jwtinfra.NewRSATokenGenerator(privateKey1, time.Hour)
	tokenStr, err := gen.Generate(uuid.New(), false)
	require.NoError(t, err)

	validator := jwtinfra.NewRSATokenValidator(&privateKey2.PublicKey)
	_, err = validator.Validate(tokenStr)

	assert.Error(t, err)
}

func TestRSATokenValidator_Validate_MalformedToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	validator := jwtinfra.NewRSATokenValidator(&privateKey.PublicKey)
	_, err = validator.Validate("not.a.jwt.token")

	assert.Error(t, err)
}
