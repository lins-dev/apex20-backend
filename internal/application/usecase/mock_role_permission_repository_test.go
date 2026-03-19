package usecase_test

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/domain/permission"
)

type mockRolePermissionRepository struct {
	createdRolePerms  []permission.RolePermission
	createRolePermErr error
}

func (m *mockRolePermissionRepository) ListRolePermissions(_ context.Context) ([]permission.RolePermission, error) {
	return m.createdRolePerms, nil
}

func (m *mockRolePermissionRepository) GetRolePermissionByID(_ context.Context, _ uuid.UUID) (permission.RolePermission, error) {
	return permission.RolePermission{}, nil
}

func (m *mockRolePermissionRepository) CreateRolePermission(_ context.Context, rp permission.RolePermission) error {
	if m.createRolePermErr != nil {
		return m.createRolePermErr
	}
	m.createdRolePerms = append(m.createdRolePerms, rp)
	return nil
}

func (m *mockRolePermissionRepository) DeleteRolePermission(_ context.Context, _ uuid.UUID, _ time.Time) (bool, error) {
	return true, nil
}
