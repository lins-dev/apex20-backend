package port

import (
	"context"

	"github.com/apex20/backend/internal/domain/permission"
)

type PermissionRepository interface {
	ExistsAny(ctx context.Context) (bool, error)
	CreatePermission(ctx context.Context, p permission.Permission) error
	CreateRolePermission(ctx context.Context, rp permission.RolePermission) error
}
