package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.PermissionUpdater = (*UpdatePermissionUseCase)(nil)

type permissionGetterUpdater interface {
	GetPermissionByID(ctx context.Context, id uuid.UUID) (permission.Permission, error)
	UpdatePermission(ctx context.Context, p permission.Permission) error
}

type UpdatePermissionUseCase struct {
	repo permissionGetterUpdater
}

func NewUpdatePermissionUseCase(repo permissionGetterUpdater) *UpdatePermissionUseCase {
	return &UpdatePermissionUseCase{repo: repo}
}

func (uc *UpdatePermissionUseCase) Execute(ctx context.Context, input port.UpdatePermissionInput) (permission.Permission, error) {
	existing, err := uc.repo.GetPermissionByID(ctx, input.ID)
	if err != nil {
		return permission.Permission{}, fmt.Errorf("getting permission: %w", err)
	}
	existing.Name = input.Name
	existing.Description = input.Description
	existing.UpdatedAt = time.Now()
	if err := uc.repo.UpdatePermission(ctx, existing); err != nil {
		return permission.Permission{}, fmt.Errorf("updating permission: %w", err)
	}
	return existing, nil
}
