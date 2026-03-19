package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

type SeedRolePermissionsUseCase struct {
	repo port.RolePermissionRepository
}

func NewSeedRolePermissionsUseCase(repo port.RolePermissionRepository) *SeedRolePermissionsUseCase {
	return &SeedRolePermissionsUseCase{repo: repo}
}

// Execute seeds all role→permission mappings using the IDs returned by SeedPermissionsUseCase.
// permissionIDs must not be nil; if nil, it means permissions were already seeded and this is a no-op.
func (uc *SeedRolePermissionsUseCase) Execute(ctx context.Context, permissionIDs map[string]uuid.UUID) error {
	if permissionIDs == nil {
		return nil
	}

	rolePermissions := map[permission.Role][]string{
		permission.RoleGM: {
			"campaign.create", "campaign.update", "campaign.delete",
			"scene.manage", "token.move.any", "token.move.own",
			"chat.send", "chat.roll", "gm.fog_control",
		},
		permission.RolePlayer: {
			"token.move.own", "chat.send", "chat.roll",
		},
		permission.RoleTrusted: {
			"token.move.own", "token.move.any",
			"chat.send", "chat.roll", "scene.manage",
		},
		permission.RoleAdmin: {
			"campaign.create", "campaign.update", "campaign.delete",
			"scene.manage", "token.move.any", "token.move.own",
			"chat.send", "chat.roll", "gm.fog_control",
			"permission.list", "permission.create", "permission.update", "permission.delete",
			"role_permission.list", "role_permission.create", "role_permission.delete",
		},
	}

	now := time.Now()

	for role, names := range rolePermissions {
		for _, name := range names {
			id, err := uuid.NewV7()
			if err != nil {
				return fmt.Errorf("generating uuid for role_permission %s/%s: %w", role, name, err)
			}
			rp := permission.RolePermission{
				ID:           id,
				Role:         role,
				PermissionID: permissionIDs[name],
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if err := uc.repo.CreateRolePermission(ctx, rp); err != nil {
				return fmt.Errorf("creating role_permission %s/%s: %w", role, name, err)
			}
		}
	}

	return nil
}
