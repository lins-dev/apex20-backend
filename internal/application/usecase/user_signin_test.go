package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/user"
)

// --- stubs ---

type stubUserByEmailGetter struct {
	u   user.User
	err error
}

func (s *stubUserByEmailGetter) GetUserByEmail(_ context.Context, _ string) (user.User, error) {
	return s.u, s.err
}

type stubPasswordVerifier struct {
	match bool
	err   error
}

func (s *stubPasswordVerifier) Hash(_ string) (string, error)           { return "", nil }
func (s *stubPasswordVerifier) Verify(_, _ string) (bool, error)        { return s.match, s.err }

// --- tests ---

func TestSignInUseCase_Execute_ReturnsTokenOnValidCredentials(t *testing.T) {
	u := user.User{
		ID:           uuid.New(),
		Email:        "hero@apex20.com",
		Name:         "Hero",
		PasswordHash: "hashed_pw",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	uc := usecase.NewSignInUseCase(
		&stubUserByEmailGetter{u: u},
		&stubPasswordVerifier{match: true},
		&stubTokenGenerator{token: "jwt.token.here"},
	)

	out, err := uc.Execute(context.Background(), port.SignInInput{
		Email:    "hero@apex20.com",
		Password: "secret123",
	})

	require.NoError(t, err)
	assert.Equal(t, u.ID, out.User.ID)
	assert.Equal(t, "jwt.token.here", out.AccessToken)
}

func TestSignInUseCase_Execute_ReturnsErrOnUnknownEmail(t *testing.T) {
	uc := usecase.NewSignInUseCase(
		&stubUserByEmailGetter{err: port.ErrNotFound},
		&stubPasswordVerifier{match: true},
		&stubTokenGenerator{token: "jwt.token.here"},
	)

	_, err := uc.Execute(context.Background(), port.SignInInput{
		Email:    "ghost@apex20.com",
		Password: "secret123",
	})

	assert.ErrorIs(t, err, port.ErrInvalidCredentials)
}

func TestSignInUseCase_Execute_ReturnsErrOnWrongPassword(t *testing.T) {
	u := user.User{ID: uuid.New(), Email: "hero@apex20.com", PasswordHash: "hashed_pw"}
	uc := usecase.NewSignInUseCase(
		&stubUserByEmailGetter{u: u},
		&stubPasswordVerifier{match: false},
		&stubTokenGenerator{token: "jwt.token.here"},
	)

	_, err := uc.Execute(context.Background(), port.SignInInput{
		Email:    "hero@apex20.com",
		Password: "wrong_password",
	})

	assert.ErrorIs(t, err, port.ErrInvalidCredentials)
}
