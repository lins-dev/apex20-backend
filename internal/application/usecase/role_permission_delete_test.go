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

type mockRolePermissionDeleter struct {
	result bool
	err    error
}

func (m *mockRolePermissionDeleter) DeleteRolePermission(_ context.Context, _ uuid.UUID, _ time.Time) (bool, error) {
	return m.result, m.err
}

func TestDeleteRolePermissionUseCase_Execute_ReturnsTrue(t *testing.T) {
	uc := usecase.NewDeleteRolePermissionUseCase(&mockRolePermissionDeleter{result: true})

	deleted, err := uc.Execute(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeleteRolePermissionUseCase_Execute_ReturnsFalseWhenNotFound(t *testing.T) {
	uc := usecase.NewDeleteRolePermissionUseCase(&mockRolePermissionDeleter{result: false})

	deleted, err := uc.Execute(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.False(t, deleted)
}

func TestDeleteRolePermissionUseCase_Execute_ReturnsError(t *testing.T) {
	uc := usecase.NewDeleteRolePermissionUseCase(&mockRolePermissionDeleter{err: errors.New("db error")})

	_, err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "deleting role permission")
}
