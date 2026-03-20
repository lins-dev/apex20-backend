package port

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/domain/permission"
)

type RolePermissionRepository interface {
	ListRolePermissions(ctx context.Context) ([]permission.RolePermission, error)
	GetRolePermissionByID(ctx context.Context, id uuid.UUID) (permission.RolePermission, error)
	CreateRolePermission(ctx context.Context, rp permission.RolePermission) error
	DeleteRolePermission(ctx context.Context, id uuid.UUID, now time.Time) (bool, error)
}
