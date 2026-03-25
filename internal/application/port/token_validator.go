package port

import "github.com/google/uuid"

// AuthClaims holds the identity claims extracted from a validated JWT.
type AuthClaims struct {
	UserID  uuid.UUID
	IsAdmin bool
}

// TokenValidator validates a JWT string and returns its claims.
type TokenValidator interface {
	Validate(tokenString string) (AuthClaims, error)
}
