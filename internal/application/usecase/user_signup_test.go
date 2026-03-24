package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/user"
)

// --- stubs ---

type stubUserCreator struct{ err error }

func (s *stubUserCreator) CreateUser(_ context.Context, _ user.User) error { return s.err }

type stubPasswordHasher struct {
	hash string
	err  error
}

func (s *stubPasswordHasher) Hash(_ string) (string, error)           { return s.hash, s.err }
func (s *stubPasswordHasher) Verify(_, _ string) (bool, error)        { return true, nil }

type stubTokenGenerator struct {
	token string
	err   error
}

func (s *stubTokenGenerator) Generate(_ uuid.UUID, _ bool) (string, error) {
	return s.token, s.err
}

// --- tests ---

func TestSignUpUseCase_Execute_CreatesUserAndReturnsToken(t *testing.T) {
	uc := usecase.NewSignUpUseCase(
		&stubUserCreator{},
		&stubPasswordHasher{hash: "hashed_pw"},
		&stubTokenGenerator{token: "jwt.token.here"},
	)

	out, err := uc.Execute(context.Background(), port.SignUpInput{
		Email:    "hero@apex20.com",
		Password: "secret123",
		Name:     "Hero",
	})

	require.NoError(t, err)
	assert.Equal(t, "hero@apex20.com", out.User.Email)
	assert.Equal(t, "Hero", out.User.Name)
	assert.NotEqual(t, uuid.Nil, out.User.ID)
	assert.Equal(t, "jwt.token.here", out.AccessToken)
}

func TestSignUpUseCase_Execute_ReturnsErrOnDuplicateEmail(t *testing.T) {
	uc := usecase.NewSignUpUseCase(
		&stubUserCreator{err: port.ErrEmailAlreadyExists},
		&stubPasswordHasher{hash: "hashed_pw"},
		&stubTokenGenerator{token: "jwt.token.here"},
	)

	_, err := uc.Execute(context.Background(), port.SignUpInput{
		Email:    "existing@apex20.com",
		Password: "secret123",
		Name:     "Hero",
	})

	assert.ErrorIs(t, err, port.ErrEmailAlreadyExists)
}

func TestSignUpUseCase_Execute_ReturnsErrOnHashFailure(t *testing.T) {
	hashErr := errors.New("hash error")
	uc := usecase.NewSignUpUseCase(
		&stubUserCreator{},
		&stubPasswordHasher{err: hashErr},
		&stubTokenGenerator{token: "jwt.token.here"},
	)

	_, err := uc.Execute(context.Background(), port.SignUpInput{
		Email:    "hero@apex20.com",
		Password: "secret123",
		Name:     "Hero",
	})

	assert.ErrorContains(t, err, "hash error")
}
