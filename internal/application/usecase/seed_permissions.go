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

// Execute seeds all platform permissions. It is idempotent — skips if any permission already exists.
// Returns the map of permission name → UUID for use by SeedRolePermissionsUseCase.
func (uc *SeedPermissionsUseCase) Execute(ctx context.Context) (map[string]uuid.UUID, error) {
	exists, err := uc.repo.ExistsAny(ctx)
	if err != nil {
		return nil, fmt.Errorf("checking existing permissions: %w", err)
	}
	if exists {
		return nil, nil
	}

	permissions := []permission.Permission{
		// Game permissions
		{Name: "campaign.create", Description: "Criar uma nova campanha"},
		{Name: "campaign.update", Description: "Editar configurações de uma campanha"},
		{Name: "campaign.delete", Description: "Deletar uma campanha"},
		{Name: "scene.manage", Description: "Gerenciar cenas (criar, editar, deletar)"},
		{Name: "token.move.any", Description: "Mover qualquer token no grid"},
		{Name: "token.move.own", Description: "Mover apenas os próprios tokens"},
		{Name: "chat.send", Description: "Enviar mensagens no chat"},
		{Name: "chat.roll", Description: "Realizar rolagens de dados no chat"},
		{Name: "gm.fog_control", Description: "Controlar névoa de guerra e visibilidade"},
		// Admin permissions
		{Name: "permission.list", Description: "Listar permissões"},
		{Name: "permission.create", Description: "Criar permissão"},
		{Name: "permission.update", Description: "Editar permissão"},
		{Name: "permission.delete", Description: "Deletar permissão"},
		{Name: "role_permission.list", Description: "Listar permissões de roles"},
		{Name: "role_permission.create", Description: "Atribuir permissão a uma role"},
		{Name: "role_permission.delete", Description: "Remover permissão de uma role"},
	}

	now := time.Now()
	ids := make(map[string]uuid.UUID, len(permissions))

	for _, p := range permissions {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("generating uuid for permission %q: %w", p.Name, err)
		}
		p.ID = id
		p.CreatedAt = now
		p.UpdatedAt = now

		if err := uc.repo.CreatePermission(ctx, p); err != nil {
			return nil, fmt.Errorf("creating permission %q: %w", p.Name, err)
		}

		ids[p.Name] = id
	}

	return ids, nil
}
