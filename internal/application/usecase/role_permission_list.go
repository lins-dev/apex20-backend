package usecase

import (
	"context"
	"fmt"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.RolePermissionLister = (*ListRolePermissionsUseCase)(nil)

type rolePermissionLister interface {
	ListRolePermissions(ctx context.Context) ([]permission.RolePermission, error)
}

type ListRolePermissionsUseCase struct {
	repo rolePermissionLister
}

func NewListRolePermissionsUseCase(repo rolePermissionLister) *ListRolePermissionsUseCase {
	return &ListRolePermissionsUseCase{repo: repo}
}

func (uc *ListRolePermissionsUseCase) Execute(ctx context.Context) ([]permission.RolePermission, error) {
	rps, err := uc.repo.ListRolePermissions(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing role permissions: %w", err)
	}
	return rps, nil
}
