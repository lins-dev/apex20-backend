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

type stubUserByIDGetter struct {
	u   user.User
	err error
}

func (s *stubUserByIDGetter) GetUserByID(_ context.Context, _ uuid.UUID) (user.User, error) {
	return s.u, s.err
}

func TestGetUserUseCase_Execute_ReturnsUser(t *testing.T) {
	id := uuid.New()
	expected := user.User{ID: id, Email: "hero@apex20.com", Name: "Hero", CreatedAt: time.Now()}
	uc := usecase.NewGetUserUseCase(&stubUserByIDGetter{u: expected})

	u, err := uc.Execute(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, expected.ID, u.ID)
	assert.Equal(t, expected.Email, u.Email)
}

func TestGetUserUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewGetUserUseCase(&stubUserByIDGetter{err: port.ErrNotFound})

	_, err := uc.Execute(context.Background(), uuid.New())

	assert.ErrorIs(t, err, port.ErrNotFound)
}
