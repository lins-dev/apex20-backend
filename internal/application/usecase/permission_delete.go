package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
)

var _ port.PermissionDeleter = (*DeletePermissionUseCase)(nil)

type permissionDeleter interface {
	DeletePermission(ctx context.Context, id uuid.UUID, now time.Time) (bool, error)
}

type DeletePermissionUseCase struct {
	repo permissionDeleter
}

func NewDeletePermissionUseCase(repo permissionDeleter) *DeletePermissionUseCase {
	return &DeletePermissionUseCase{repo: repo}
}

func (uc *DeletePermissionUseCase) Execute(ctx context.Context, id uuid.UUID) (bool, error) {
	deleted, err := uc.repo.DeletePermission(ctx, id, time.Now())
	if err != nil {
		return false, fmt.Errorf("deleting permission: %w", err)
	}
	return deleted, nil
}
