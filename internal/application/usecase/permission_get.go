package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.PermissionGetter = (*GetPermissionUseCase)(nil)

type permissionGetter interface {
	GetPermissionByID(ctx context.Context, id uuid.UUID) (permission.Permission, error)
}

type GetPermissionUseCase struct {
	repo permissionGetter
}

func NewGetPermissionUseCase(repo permissionGetter) *GetPermissionUseCase {
	return &GetPermissionUseCase{repo: repo}
}

func (uc *GetPermissionUseCase) Execute(ctx context.Context, id uuid.UUID) (permission.Permission, error) {
	p, err := uc.repo.GetPermissionByID(ctx, id)
	if err != nil {
		return permission.Permission{}, fmt.Errorf("getting permission: %w", err)
	}
	return p, nil
}
