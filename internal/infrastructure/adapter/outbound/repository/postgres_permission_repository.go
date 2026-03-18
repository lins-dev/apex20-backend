package repository

import (
	"context"
	"database/sql"

	repositorygen "github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository/gen"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.PermissionRepository = (*PostgresPermissionRepository)(nil)

type PostgresPermissionRepository struct {
	queries *repositorygen.Queries
}

func NewPostgresPermissionRepository(db *sql.DB) *PostgresPermissionRepository {
	return &PostgresPermissionRepository{queries: repositorygen.New(db)}
}

func (r *PostgresPermissionRepository) ExistsAny(ctx context.Context) (bool, error) {
	return r.queries.ExistsAnyPermission(ctx)
}

func (r *PostgresPermissionRepository) CreatePermission(ctx context.Context, p permission.Permission) error {
	return r.queries.CreatePermission(ctx, repositorygen.CreatePermissionParams{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	})
}

func (r *PostgresPermissionRepository) CreateRolePermission(ctx context.Context, rp permission.RolePermission) error {
	return r.queries.CreateRolePermission(ctx, repositorygen.CreateRolePermissionParams{
		ID:           rp.ID,
		Role:         repositorygen.UserRole(rp.Role),
		PermissionID: rp.PermissionID,
		CreatedAt:    rp.CreatedAt,
		UpdatedAt:    rp.UpdatedAt,
	})
}
