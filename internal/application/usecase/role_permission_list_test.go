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
	"github.com/apex20/backend/internal/domain/permission"
)

type mockRolePermissionLister struct {
	result []permission.RolePermission
	err    error
}

func (m *mockRolePermissionLister) ListRolePermissions(_ context.Context) ([]permission.RolePermission, error) {
	return m.result, m.err
}

func TestListRolePermissionsUseCase_Execute_ReturnsList(t *testing.T) {
	rps := []permission.RolePermission{
		{ID: uuid.New(), Role: permission.RoleGM, PermissionID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Role: permission.RolePlayer, PermissionID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	uc := usecase.NewListRolePermissionsUseCase(&mockRolePermissionLister{result: rps})

	got, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestListRolePermissionsUseCase_Execute_ReturnsError(t *testing.T) {
	uc := usecase.NewListRolePermissionsUseCase(&mockRolePermissionLister{err: errors.New("db error")})

	got, err := uc.Execute(context.Background())

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "listing role permissions")
}
