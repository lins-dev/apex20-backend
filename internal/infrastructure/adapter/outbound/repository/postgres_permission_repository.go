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

func (r *PostgresPermissionRepository) ListPermissions(ctx context.Context) ([]permission.Permission, error) {
	rows, err := r.queries.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]permission.Permission, len(rows))
	for i, row := range rows {
		result[i] = toPermissionDomain(row)
	}
	return result, nil
}

func (r *PostgresPermissionRepository) GetPermissionByID(ctx context.Context, id uuid.UUID) (permission.Permission, error) {
	row, err := r.queries.GetPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return permission.Permission{}, port.ErrNotFound
		}
		return permission.Permission{}, err
	}
	return toPermissionDomain(row), nil
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

func (r *PostgresPermissionRepository) UpdatePermission(ctx context.Context, p permission.Permission) error {
	return r.queries.UpdatePermission(ctx, repositorygen.UpdatePermissionParams{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		UpdatedAt:   p.UpdatedAt,
	})
}

func (r *PostgresPermissionRepository) DeletePermission(ctx context.Context, id uuid.UUID, now time.Time) (bool, error) {
	rows, err := r.queries.SoftDeletePermission(ctx, repositorygen.SoftDeletePermissionParams{
		ID:        id,
		DeletedAt: sql.NullTime{Time: now, Valid: true},
	})
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func toPermissionDomain(row repositorygen.Permission) permission.Permission {
	return permission.Permission{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
