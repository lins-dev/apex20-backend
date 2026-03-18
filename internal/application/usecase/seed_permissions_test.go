package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/permission"
)

type mockPermissionRepository struct {
	existsAny          bool
	existsAnyErr       error
	createdPermissions []permission.Permission
	createPermErr      error
	createdRolePerms   []permission.RolePermission
	createRolePermErr  error
}

func (m *mockPermissionRepository) ExistsAny(_ context.Context) (bool, error) {
	return m.existsAny, m.existsAnyErr
}

func (m *mockPermissionRepository) CreatePermission(_ context.Context, p permission.Permission) error {
	if m.createPermErr != nil {
		return m.createPermErr
	}
	m.createdPermissions = append(m.createdPermissions, p)
	return nil
}

func (m *mockPermissionRepository) CreateRolePermission(_ context.Context, rp permission.RolePermission) error {
	if m.createRolePermErr != nil {
		return m.createRolePermErr
	}
	m.createdRolePerms = append(m.createdRolePerms, rp)
	return nil
}

func TestSeedPermissionsUseCase_Execute_SkipsWhenAlreadySeeded(t *testing.T) {
	repo := &mockPermissionRepository{existsAny: true}
	uc := usecase.NewSeedPermissionsUseCase(repo)

	err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Empty(t, repo.createdPermissions)
	assert.Empty(t, repo.createdRolePerms)
}

func TestSeedPermissionsUseCase_Execute_CreatesPermissionsAndRolePermissions(t *testing.T) {
	repo := &mockPermissionRepository{existsAny: false}
	uc := usecase.NewSeedPermissionsUseCase(repo)

	err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Len(t, repo.createdPermissions, 9)
	assert.NotEmpty(t, repo.createdRolePerms)

	// Verifica que todos os IDs são únicos (UUIDv7)
	ids := make(map[string]struct{}, len(repo.createdPermissions))
	for _, p := range repo.createdPermissions {
		id := p.ID.String()
		assert.NotEqual(t, "00000000-0000-0000-0000-000000000000", id)
		ids[id] = struct{}{}
	}
	assert.Len(t, ids, 9, "todos os IDs de permission devem ser únicos")
}

func TestSeedPermissionsUseCase_Execute_ReturnsErrorOnExistsAnyFailure(t *testing.T) {
	repo := &mockPermissionRepository{existsAnyErr: errors.New("db error")}
	uc := usecase.NewSeedPermissionsUseCase(repo)

	err := uc.Execute(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "checking existing permissions")
}

func TestSeedPermissionsUseCase_Execute_ReturnsErrorOnCreatePermissionFailure(t *testing.T) {
	repo := &mockPermissionRepository{createPermErr: errors.New("insert failed")}
	uc := usecase.NewSeedPermissionsUseCase(repo)

	err := uc.Execute(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "creating permission")
}

func TestSeedPermissionsUseCase_Execute_ReturnsErrorOnCreateRolePermissionFailure(t *testing.T) {
	repo := &mockPermissionRepository{createRolePermErr: errors.New("insert failed")}
	uc := usecase.NewSeedPermissionsUseCase(repo)

	err := uc.Execute(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "creating role_permission")
}
