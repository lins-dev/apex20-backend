package jwt

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
)

// RSATokenValidator implements port.TokenValidator using RS256.
type RSATokenValidator struct {
	publicKey *rsa.PublicKey
}

func NewRSATokenValidator(publicKey *rsa.PublicKey) *RSATokenValidator {
	return &RSATokenValidator{publicKey: publicKey}
}

func (v *RSATokenValidator) Validate(tokenString string) (port.AuthClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return v.publicKey, nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return port.AuthClaims{}, fmt.Errorf("validating token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return port.AuthClaims{}, fmt.Errorf("invalid claims type")
	}

	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return port.AuthClaims{}, fmt.Errorf("missing sub claim")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return port.AuthClaims{}, fmt.Errorf("invalid sub claim: %w", err)
	}

	isAdmin, _ := claims["is_admin"].(bool)

	return port.AuthClaims{UserID: userID, IsAdmin: isAdmin}, nil
}
