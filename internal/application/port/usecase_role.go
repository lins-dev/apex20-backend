package port

import (
	"context"

	"github.com/apex20/backend/internal/domain/permission"
)

// Driving port — used by inbound adapters (HTTP handlers) to call use cases.

type RoleLister interface {
	Execute(ctx context.Context) ([]permission.Role, error)
}
