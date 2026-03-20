package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/usecase"
)

func TestSeedPermissionsUseCase_Execute_SkipsWhenAlreadySeeded(t *testing.T) {
	repo := &mockPermissionRepository{existsAny: true}
	uc := usecase.NewSeedPermissionsUseCase(repo)

	ids, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Nil(t, ids)
	assert.Empty(t, repo.createdPermissions)
}

func TestSeedPermissionsUseCase_Execute_CreatesAllPermissions(t *testing.T) {
	repo := &mockPermissionRepository{existsAny: false}
	uc := usecase.NewSeedPermissionsUseCase(repo)

	ids, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Len(t, repo.createdPermissions, 16)
	assert.Len(t, ids, 16)

	// Todos os IDs devem ser UUIDv7 não-zero e únicos
	seen := make(map[string]struct{}, 16)
	for _, p := range repo.createdPermissions {
		id := p.ID.String()
		assert.NotEqual(t, "00000000-0000-0000-0000-000000000000", id)
		seen[id] = struct{}{}
	}
	assert.Len(t, seen, 16, "todos os IDs de permission devem ser únicos")
}

func TestSeedPermissionsUseCase_Execute_ReturnsErrorOnExistsAnyFailure(t *testing.T) {
	repo := &mockPermissionRepository{existsAnyErr: errors.New("db error")}
	uc := usecase.NewSeedPermissionsUseCase(repo)

	ids, err := uc.Execute(context.Background())

	require.Error(t, err)
	assert.Nil(t, ids)
	assert.Contains(t, err.Error(), "checking existing permissions")
}

func TestSeedPermissionsUseCase_Execute_ReturnsErrorOnCreatePermissionFailure(t *testing.T) {
	repo := &mockPermissionRepository{createPermErr: errors.New("insert failed")}
	uc := usecase.NewSeedPermissionsUseCase(repo)

	ids, err := uc.Execute(context.Background())

	require.Error(t, err)
	assert.Nil(t, ids)
	assert.Contains(t, err.Error(), "creating permission")
}
