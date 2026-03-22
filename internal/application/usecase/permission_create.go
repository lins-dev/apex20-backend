package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.PermissionCreator = (*CreatePermissionUseCase)(nil)

type permissionCreator interface {
	CreatePermission(ctx context.Context, p permission.Permission) error
}

type CreatePermissionUseCase struct {
	repo permissionCreator
}

func NewCreatePermissionUseCase(repo permissionCreator) *CreatePermissionUseCase {
	return &CreatePermissionUseCase{repo: repo}
}

func (uc *CreatePermissionUseCase) Execute(ctx context.Context, input port.CreatePermissionInput) (permission.Permission, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return permission.Permission{}, fmt.Errorf("generating id: %w", err)
	}
	now := time.Now()
	p := permission.Permission{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.repo.CreatePermission(ctx, p); err != nil {
		return permission.Permission{}, fmt.Errorf("creating permission: %w", err)
	}
	return p, nil
}
