package usecase_test

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/domain/permission"
)

type mockPermissionRepository struct {
	existsAny          bool
	existsAnyErr       error
	createdPermissions []permission.Permission
	createPermErr      error
}

func (m *mockPermissionRepository) ExistsAny(_ context.Context) (bool, error) {
	return m.existsAny, m.existsAnyErr
}

func (m *mockPermissionRepository) ListPermissions(_ context.Context) ([]permission.Permission, error) {
	return m.createdPermissions, nil
}

func (m *mockPermissionRepository) GetPermissionByID(_ context.Context, _ uuid.UUID) (permission.Permission, error) {
	return permission.Permission{}, nil
}

func (m *mockPermissionRepository) CreatePermission(_ context.Context, p permission.Permission) error {
	if m.createPermErr != nil {
		return m.createPermErr
	}
	m.createdPermissions = append(m.createdPermissions, p)
	return nil
}

func (m *mockPermissionRepository) UpdatePermission(_ context.Context, _ permission.Permission) error {
	return nil
}

func (m *mockPermissionRepository) DeletePermission(_ context.Context, _ uuid.UUID, _ time.Time) (bool, error) {
	return true, nil
}
