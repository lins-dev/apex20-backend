package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/permission"
)

func buildPermissionIDs() map[string]uuid.UUID {
	names := []string{
		"campaign.create", "campaign.update", "campaign.delete",
		"scene.manage", "token.move.any", "token.move.own",
		"chat.send", "chat.roll", "gm.fog_control",
	}
	ids := make(map[string]uuid.UUID, len(names))
	for _, n := range names {
		ids[n] = uuid.New()
	}
	return ids
}

func TestSeedRolePermissionsUseCase_Execute_SkipsWhenIDsNil(t *testing.T) {
	repo := &mockRolePermissionRepository{}
	uc := usecase.NewSeedRolePermissionsUseCase(repo)

	err := uc.Execute(context.Background(), nil)

	require.NoError(t, err)
	assert.Empty(t, repo.createdRolePerms)
}

func TestSeedRolePermissionsUseCase_Execute_CreatesAllRolePermissions(t *testing.T) {
	repo := &mockRolePermissionRepository{}
	uc := usecase.NewSeedRolePermissionsUseCase(repo)

	err := uc.Execute(context.Background(), buildPermissionIDs())

	require.NoError(t, err)
	assert.NotEmpty(t, repo.createdRolePerms)

	// Verifica que todas as 3 roles de campanha foram mapeadas
	roles := make(map[permission.Role]struct{})
	for _, rp := range repo.createdRolePerms {
		roles[rp.Role] = struct{}{}
	}
	assert.Contains(t, roles, permission.RoleGM)
	assert.Contains(t, roles, permission.RolePlayer)
	assert.Contains(t, roles, permission.RoleTrusted)

	// GM deve ter mais permissões que player
	var gmCount, playerCount int
	for _, rp := range repo.createdRolePerms {
		switch rp.Role {
		case permission.RoleGM:
			gmCount++
		case permission.RolePlayer:
			playerCount++
		}
	}
	assert.Greater(t, gmCount, playerCount)
}

func TestSeedRolePermissionsUseCase_Execute_ReturnsErrorOnCreateFailure(t *testing.T) {
	repo := &mockRolePermissionRepository{createRolePermErr: errors.New("insert failed")}
	uc := usecase.NewSeedRolePermissionsUseCase(repo)

	err := uc.Execute(context.Background(), buildPermissionIDs())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "creating role_permission")
}
