package repository_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/domain/permission"
	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository"

	"github.com/google/uuid"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	dsn := os.Getenv("DATABASE_URL")
	require.NotEmpty(t, dsn, "DATABASE_URL is required for integration tests")

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping(), "failed to connect to database")
	t.Cleanup(func() { db.Close() })
	return db
}

func TestPostgresPermissionRepository_ExistsAny_ReturnsFalseWhenEmpty(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	_, err := db.ExecContext(ctx, "DELETE FROM role_permissions")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "DELETE FROM permissions")
	require.NoError(t, err)

	repo := repository.NewPostgresPermissionRepository(db)
	exists, err := repo.ExistsAny(ctx)

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestPostgresPermissionRepository_CreatePermission_AndExistsAny(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	_, err := db.ExecContext(ctx, "DELETE FROM role_permissions")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "DELETE FROM permissions")
	require.NoError(t, err)

	repo := repository.NewPostgresPermissionRepository(db)

	id, err := uuid.NewV7()
	require.NoError(t, err)

	p := permission.Permission{
		ID:          id,
		Name:        "test.permission",
		Description: "Permissão de teste",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.CreatePermission(ctx, p)
	require.NoError(t, err)

	exists, err := repo.ExistsAny(ctx)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestPostgresPermissionRepository_CreateRolePermission(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	_, err := db.ExecContext(ctx, "DELETE FROM role_permissions")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "DELETE FROM permissions")
	require.NoError(t, err)

	repo := repository.NewPostgresPermissionRepository(db)
	now := time.Now()

	permID, err := uuid.NewV7()
	require.NoError(t, err)

	err = repo.CreatePermission(ctx, permission.Permission{
		ID:          permID,
		Name:        "chat.send",
		Description: "Enviar mensagens",
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	require.NoError(t, err)

	rpID, err := uuid.NewV7()
	require.NoError(t, err)

	rp := permission.RolePermission{
		ID:           rpID,
		Role:         permission.RolePlayer,
		PermissionID: permID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = repo.CreateRolePermission(ctx, rp)
	require.NoError(t, err)

	var count int
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM role_permissions WHERE id = $1", rpID,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
