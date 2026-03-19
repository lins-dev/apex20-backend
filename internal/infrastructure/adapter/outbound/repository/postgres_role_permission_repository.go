package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
	repositorygen "github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository/gen"
)

var _ port.RolePermissionRepository = (*PostgresRolePermissionRepository)(nil)

type PostgresRolePermissionRepository struct {
	queries *repositorygen.Queries
}

func NewPostgresRolePermissionRepository(db *sql.DB) *PostgresRolePermissionRepository {
	return &PostgresRolePermissionRepository{queries: repositorygen.New(db)}
}

func (r *PostgresRolePermissionRepository) ListRolePermissions(ctx context.Context) ([]permission.RolePermission, error) {
	rows, err := r.queries.ListRolePermissions(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]permission.RolePermission, len(rows))
	for i, row := range rows {
		result[i] = toRolePermissionDomain(row)
	}
	return result, nil
}

func (r *PostgresRolePermissionRepository) GetRolePermissionByID(ctx context.Context, id uuid.UUID) (permission.RolePermission, error) {
	row, err := r.queries.GetRolePermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return permission.RolePermission{}, port.ErrNotFound
		}
		return permission.RolePermission{}, err
	}
	return toRolePermissionDomain(row), nil
}

func (r *PostgresRolePermissionRepository) CreateRolePermission(ctx context.Context, rp permission.RolePermission) error {
	return r.queries.CreateRolePermission(ctx, repositorygen.CreateRolePermissionParams{
		ID:           rp.ID,
		Role:         rp.Role,
		PermissionID: rp.PermissionID,
		CreatedAt:    rp.CreatedAt,
		UpdatedAt:    rp.UpdatedAt,
	})
}

func (r *PostgresRolePermissionRepository) DeleteRolePermission(ctx context.Context, id uuid.UUID, now time.Time) (bool, error) {
	rows, err := r.queries.SoftDeleteRolePermission(ctx, repositorygen.SoftDeleteRolePermissionParams{
		ID:        id,
		DeletedAt: sql.NullTime{Time: now, Valid: true},
	})
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func toRolePermissionDomain(row repositorygen.RolePermission) permission.RolePermission {
	return permission.RolePermission{
		ID:           row.ID,
		Role:         row.Role,
		PermissionID: row.PermissionID,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
