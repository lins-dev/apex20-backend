package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// RSATokenGenerator implements port.TokenGenerator using RS256.
type RSATokenGenerator struct {
	privateKey *rsa.PrivateKey
	expiry     time.Duration
}

func NewRSATokenGenerator(privateKey *rsa.PrivateKey, expiry time.Duration) *RSATokenGenerator {
	return &RSATokenGenerator{privateKey: privateKey, expiry: expiry}
}

func (g *RSATokenGenerator) Generate(userID uuid.UUID, isAdmin bool) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      userID.String(),
		"is_admin": isAdmin,
		"iat":      now.Unix(),
		"exp":      now.Add(g.expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(g.privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return signed, nil
}
