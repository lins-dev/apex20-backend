package usecase

import (
	"context"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.RoleLister = (*ListRolesUseCase)(nil)

type ListRolesUseCase struct{}

func NewListRolesUseCase() *ListRolesUseCase {
	return &ListRolesUseCase{}
}

func (uc *ListRolesUseCase) Execute(_ context.Context) ([]permission.Role, error) {
	return []permission.Role{
		permission.RoleGM,
		permission.RolePlayer,
		permission.RoleTrusted,
	}, nil
}
