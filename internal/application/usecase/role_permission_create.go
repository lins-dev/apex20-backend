package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.RolePermissionCreator = (*CreateRolePermissionUseCase)(nil)

type rolePermissionCreator interface {
	CreateRolePermission(ctx context.Context, rp permission.RolePermission) error
}

type CreateRolePermissionUseCase struct {
	repo rolePermissionCreator
}

func NewCreateRolePermissionUseCase(repo rolePermissionCreator) *CreateRolePermissionUseCase {
	return &CreateRolePermissionUseCase{repo: repo}
}

func (uc *CreateRolePermissionUseCase) Execute(ctx context.Context, input port.CreateRolePermissionInput) (permission.RolePermission, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return permission.RolePermission{}, fmt.Errorf("generating id: %w", err)
	}
	now := time.Now()
	rp := permission.RolePermission{
		ID:           id,
		Role:         input.Role,
		PermissionID: input.PermissionID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := uc.repo.CreateRolePermission(ctx, rp); err != nil {
		return permission.RolePermission{}, fmt.Errorf("creating role permission: %w", err)
	}
	return rp, nil
}
