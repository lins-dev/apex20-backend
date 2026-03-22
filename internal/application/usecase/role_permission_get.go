package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.RolePermissionGetter = (*GetRolePermissionUseCase)(nil)

type rolePermissionGetter interface {
	GetRolePermissionByID(ctx context.Context, id uuid.UUID) (permission.RolePermission, error)
}

type GetRolePermissionUseCase struct {
	repo rolePermissionGetter
}

func NewGetRolePermissionUseCase(repo rolePermissionGetter) *GetRolePermissionUseCase {
	return &GetRolePermissionUseCase{repo: repo}
}

func (uc *GetRolePermissionUseCase) Execute(ctx context.Context, id uuid.UUID) (permission.RolePermission, error) {
	rp, err := uc.repo.GetRolePermissionByID(ctx, id)
	if err != nil {
		return permission.RolePermission{}, fmt.Errorf("getting role permission: %w", err)
	}
	return rp, nil
}
