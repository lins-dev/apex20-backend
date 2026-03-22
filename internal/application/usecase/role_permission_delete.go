package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
)

var _ port.RolePermissionDeleter = (*DeleteRolePermissionUseCase)(nil)

type rolePermissionDeleter interface {
	DeleteRolePermission(ctx context.Context, id uuid.UUID, now time.Time) (bool, error)
}

type DeleteRolePermissionUseCase struct {
	repo rolePermissionDeleter
}

func NewDeleteRolePermissionUseCase(repo rolePermissionDeleter) *DeleteRolePermissionUseCase {
	return &DeleteRolePermissionUseCase{repo: repo}
}

func (uc *DeleteRolePermissionUseCase) Execute(ctx context.Context, id uuid.UUID) (bool, error) {
	deleted, err := uc.repo.DeleteRolePermission(ctx, id, time.Now())
	if err != nil {
		return false, fmt.Errorf("deleting role permission: %w", err)
	}
	return deleted, nil
}
