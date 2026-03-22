package usecase

import (
	"context"
	"fmt"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.PermissionLister = (*ListPermissionsUseCase)(nil)

type permissionLister interface {
	ListPermissions(ctx context.Context) ([]permission.Permission, error)
}

type ListPermissionsUseCase struct {
	repo permissionLister
}

func NewListPermissionsUseCase(repo permissionLister) *ListPermissionsUseCase {
	return &ListPermissionsUseCase{repo: repo}
}

func (uc *ListPermissionsUseCase) Execute(ctx context.Context) ([]permission.Permission, error) {
	perms, err := uc.repo.ListPermissions(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing permissions: %w", err)
	}
	return perms, nil
}
