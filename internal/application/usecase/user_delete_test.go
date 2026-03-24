package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
)

type stubUserSoftDeleter struct{ err error }

func (s *stubUserSoftDeleter) DeleteUser(_ context.Context, _ uuid.UUID) error { return s.err }

func TestDeleteUserUseCase_Execute_Succeeds(t *testing.T) {
	uc := usecase.NewDeleteUserUseCase(&stubUserSoftDeleter{})

	err := uc.Execute(context.Background(), uuid.New())

	require.NoError(t, err)
}

func TestDeleteUserUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewDeleteUserUseCase(&stubUserSoftDeleter{err: port.ErrNotFound})

	err := uc.Execute(context.Background(), uuid.New())

	assert.ErrorIs(t, err, port.ErrNotFound)
}
