package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository"
)

func TestPostgresRolePermissionRepository_CreateAndList(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	permRepo := repository.NewPostgresPermissionRepository(db)
	rolePermRepo := repository.NewPostgresRolePermissionRepository(db)
	ctx := context.Background()

	perm := createTestPermission(t, permRepo, "chat.send")

	rpID, err := uuid.NewV7()
	require.NoError(t, err)
	now := time.Now()

	rp := permission.RolePermission{
		ID:           rpID,
		Role:         permission.RolePlayer,
		PermissionID: perm.ID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, rolePermRepo.CreateRolePermission(ctx, rp))

	list, err := rolePermRepo.ListRolePermissions(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, permission.RolePlayer, list[0].Role)
	assert.Equal(t, perm.ID, list[0].PermissionID)
}

func TestPostgresRolePermissionRepository_GetByID(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	permRepo := repository.NewPostgresPermissionRepository(db)
	rolePermRepo := repository.NewPostgresRolePermissionRepository(db)
	ctx := context.Background()

	perm := createTestPermission(t, permRepo, "token.move.own")

	rpID, err := uuid.NewV7()
	require.NoError(t, err)
	now := time.Now()

	rp := permission.RolePermission{
		ID:           rpID,
		Role:         permission.RoleGM,
		PermissionID: perm.ID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, rolePermRepo.CreateRolePermission(ctx, rp))

	got, err := rolePermRepo.GetRolePermissionByID(ctx, rpID)
	require.NoError(t, err)
	assert.Equal(t, rpID, got.ID)
	assert.Equal(t, permission.RoleGM, got.Role)
}

func TestPostgresRolePermissionRepository_GetByID_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	rolePermRepo := repository.NewPostgresRolePermissionRepository(db)
	_, err := rolePermRepo.GetRolePermissionByID(context.Background(), uuid.New())

	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestPostgresRolePermissionRepository_Delete(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	permRepo := repository.NewPostgresPermissionRepository(db)
	rolePermRepo := repository.NewPostgresRolePermissionRepository(db)
	ctx := context.Background()

	perm := createTestPermission(t, permRepo, "scene.manage")

	rpID, err := uuid.NewV7()
	require.NoError(t, err)
	now := time.Now()

	require.NoError(t, rolePermRepo.CreateRolePermission(ctx, permission.RolePermission{
		ID:           rpID,
		Role:         permission.RoleTrusted,
		PermissionID: perm.ID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}))

	deleted, err := rolePermRepo.DeleteRolePermission(ctx, rpID, time.Now())
	require.NoError(t, err)
	assert.True(t, deleted)

	_, err = rolePermRepo.GetRolePermissionByID(ctx, rpID)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestPostgresRolePermissionRepository_Delete_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	rolePermRepo := repository.NewPostgresRolePermissionRepository(db)
	deleted, err := rolePermRepo.DeleteRolePermission(context.Background(), uuid.New(), time.Now())

	require.NoError(t, err)
	assert.False(t, deleted)
}
