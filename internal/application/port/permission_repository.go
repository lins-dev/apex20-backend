package port

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/domain/permission"
)

type PermissionRepository interface {
	ExistsAny(ctx context.Context) (bool, error)
	ListPermissions(ctx context.Context) ([]permission.Permission, error)
	GetPermissionByID(ctx context.Context, id uuid.UUID) (permission.Permission, error)
	CreatePermission(ctx context.Context, p permission.Permission) error
	UpdatePermission(ctx context.Context, p permission.Permission) error
	DeletePermission(ctx context.Context, id uuid.UUID, now time.Time) (bool, error)
}
