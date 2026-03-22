package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/usecase"
)

type mockPermissionDeleter struct {
	result bool
	err    error
}

func (m *mockPermissionDeleter) DeletePermission(_ context.Context, _ uuid.UUID, _ time.Time) (bool, error) {
	return m.result, m.err
}

func TestDeletePermissionUseCase_Execute_ReturnsTrue(t *testing.T) {
	repo := &mockPermissionDeleter{result: true}
	uc := usecase.NewDeletePermissionUseCase(repo)

	deleted, err := uc.Execute(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeletePermissionUseCase_Execute_ReturnsFalseWhenNotFound(t *testing.T) {
	repo := &mockPermissionDeleter{result: false}
	uc := usecase.NewDeletePermissionUseCase(repo)

	deleted, err := uc.Execute(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.False(t, deleted)
}

func TestDeletePermissionUseCase_Execute_ReturnsError(t *testing.T) {
	repo := &mockPermissionDeleter{err: errors.New("db error")}
	uc := usecase.NewDeletePermissionUseCase(repo)

	_, err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "deleting permission")
}
