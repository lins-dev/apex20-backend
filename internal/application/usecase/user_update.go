package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/user"
)

var _ port.UserUpdater = (*UpdateUserUseCase)(nil)

type userUpdater interface {
	UpdateUser(ctx context.Context, id uuid.UUID, name string, nick *string) (user.User, error)
}

type UpdateUserUseCase struct {
	repo userUpdater
}

func NewUpdateUserUseCase(repo userUpdater) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input port.UpdateUserInput) (user.User, error) {
	return uc.repo.UpdateUser(ctx, input.ID, input.Name, input.Nick)
}
