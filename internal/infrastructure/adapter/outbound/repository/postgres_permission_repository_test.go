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

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository"
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

func cleanDB(t *testing.T, db *sql.DB) {
	t.Helper()
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "DELETE FROM role_permissions")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "DELETE FROM permissions")
	require.NoError(t, err)
}

func createTestPermission(t *testing.T, repo *repository.PostgresPermissionRepository, name string) permission.Permission {
	t.Helper()
	id, err := uuid.NewV7()
	require.NoError(t, err)
	now := time.Now()
	p := permission.Permission{
		ID:          id,
		Name:        name,
		Description: "Permissão de teste",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, repo.CreatePermission(context.Background(), p))
	return p
}

func TestPostgresPermissionRepository_ExistsAny_ReturnsFalseWhenEmpty(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresPermissionRepository(db)
	exists, err := repo.ExistsAny(context.Background())

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestPostgresPermissionRepository_CreateAndExistsAny(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresPermissionRepository(db)
	createTestPermission(t, repo, "test.permission")

	exists, err := repo.ExistsAny(context.Background())
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestPostgresPermissionRepository_ListPermissions(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresPermissionRepository(db)
	createTestPermission(t, repo, "perm.a")
	createTestPermission(t, repo, "perm.b")

	list, err := repo.ListPermissions(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestPostgresPermissionRepository_GetPermissionByID(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresPermissionRepository(db)
	created := createTestPermission(t, repo, "get.test")

	got, err := repo.GetPermissionByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, created.Name, got.Name)
}

func TestPostgresPermissionRepository_GetPermissionByID_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresPermissionRepository(db)
	_, err := repo.GetPermissionByID(context.Background(), uuid.New())

	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestPostgresPermissionRepository_UpdatePermission(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresPermissionRepository(db)
	created := createTestPermission(t, repo, "update.test")

	updated := created
	updated.Name = "update.test.renamed"
	updated.Description = "Nova descrição"
	updated.UpdatedAt = time.Now()

	require.NoError(t, repo.UpdatePermission(context.Background(), updated))

	got, err := repo.GetPermissionByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, "update.test.renamed", got.Name)
	assert.Equal(t, "Nova descrição", got.Description)
}

func TestPostgresPermissionRepository_DeletePermission(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresPermissionRepository(db)
	created := createTestPermission(t, repo, "delete.test")

	deleted, err := repo.DeletePermission(context.Background(), created.ID, time.Now())
	require.NoError(t, err)
	assert.True(t, deleted)

	_, err = repo.GetPermissionByID(context.Background(), created.ID)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestPostgresPermissionRepository_DeletePermission_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresPermissionRepository(db)
	deleted, err := repo.DeletePermission(context.Background(), uuid.New(), time.Now())

	require.NoError(t, err)
	assert.False(t, deleted)
}
