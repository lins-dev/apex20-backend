package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/permission"
)

type mockRolePermissionGetter struct {
	result permission.RolePermission
	err    error
}

func (m *mockRolePermissionGetter) GetRolePermissionByID(_ context.Context, _ uuid.UUID) (permission.RolePermission, error) {
	return m.result, m.err
}

func TestGetRolePermissionUseCase_Execute_ReturnsRolePermission(t *testing.T) {
	id := uuid.New()
	expected := permission.RolePermission{ID: id, Role: permission.RoleGM, PermissionID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	uc := usecase.NewGetRolePermissionUseCase(&mockRolePermissionGetter{result: expected})

	got, err := uc.Execute(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Role, got.Role)
}

func TestGetRolePermissionUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewGetRolePermissionUseCase(&mockRolePermissionGetter{err: port.ErrNotFound})

	_, err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.True(t, errors.Is(err, port.ErrNotFound))
}

func TestGetRolePermissionUseCase_Execute_ReturnsRepoError(t *testing.T) {
	uc := usecase.NewGetRolePermissionUseCase(&mockRolePermissionGetter{err: errors.New("db error")})

	_, err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting role permission")
}
