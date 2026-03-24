package port

import "github.com/google/uuid"

// PasswordHasher hashes and verifies passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

// TokenGenerator generates a signed JWT for a given user.
type TokenGenerator interface {
	Generate(userID uuid.UUID, isAdmin bool) (string, error)
}
