package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

type SeedPermissionsUseCase struct {
	repo port.PermissionRepository
}

func NewSeedPermissionsUseCase(repo port.PermissionRepository) *SeedPermissionsUseCase {
	return &SeedPermissionsUseCase{repo: repo}
}

func (uc *SeedPermissionsUseCase) Execute(ctx context.Context) error {
	exists, err := uc.repo.ExistsAny(ctx)
	if err != nil {
		return fmt.Errorf("checking existing permissions: %w", err)
	}
	if exists {
		return nil
	}

	permissions := []permission.Permission{
		{Name: "campaign.create", Description: "Criar uma nova campanha"},
		{Name: "campaign.update", Description: "Editar configurações de uma campanha"},
		{Name: "campaign.delete", Description: "Deletar uma campanha"},
		{Name: "scene.manage", Description: "Gerenciar cenas (criar, editar, deletar)"},
		{Name: "token.move.any", Description: "Mover qualquer token no grid"},
		{Name: "token.move.own", Description: "Mover apenas os próprios tokens"},
		{Name: "chat.send", Description: "Enviar mensagens no chat"},
		{Name: "chat.roll", Description: "Realizar rolagens de dados no chat"},
		{Name: "gm.fog_control", Description: "Controlar névoa de guerra e visibilidade"},
	}

	now := time.Now()
	permissionIDs := make(map[string]uuid.UUID, len(permissions))

	for _, p := range permissions {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generating uuid for permission %q: %w", p.Name, err)
		}
		p.ID = id
		p.CreatedAt = now
		p.UpdatedAt = now

		if err := uc.repo.CreatePermission(ctx, p); err != nil {
			return fmt.Errorf("creating permission %q: %w", p.Name, err)
		}

		permissionIDs[p.Name] = id
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
	}

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
