package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/domain/permission"
)

// Input types

type CreatePermissionInput struct {
	Name        string
	Description string
}

type UpdatePermissionInput struct {
	ID          uuid.UUID
	Name        string
	Description string
}

// Driving ports — used by inbound adapters (HTTP handlers) to call use cases.

type PermissionLister interface {
	Execute(ctx context.Context) ([]permission.Permission, error)
}

type PermissionGetter interface {
	Execute(ctx context.Context, id uuid.UUID) (permission.Permission, error)
}

type PermissionCreator interface {
	Execute(ctx context.Context, input CreatePermissionInput) (permission.Permission, error)
}

type PermissionUpdater interface {
	Execute(ctx context.Context, input UpdatePermissionInput) (permission.Permission, error)
}

type PermissionDeleter interface {
	Execute(ctx context.Context, id uuid.UUID) (bool, error)
}
