package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/domain/permission"
)

// Input types

type CreateRolePermissionInput struct {
	Role         permission.Role
	PermissionID uuid.UUID
}

// Driving ports — used by inbound adapters (HTTP handlers) to call use cases.

type RolePermissionLister interface {
	Execute(ctx context.Context) ([]permission.RolePermission, error)
}

type RolePermissionGetter interface {
	Execute(ctx context.Context, id uuid.UUID) (permission.RolePermission, error)
}

type RolePermissionCreator interface {
	Execute(ctx context.Context, input CreateRolePermissionInput) (permission.RolePermission, error)
}

type RolePermissionDeleter interface {
	Execute(ctx context.Context, id uuid.UUID) (bool, error)
}
