package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
)

var _ port.UserDeleter = (*DeleteUserUseCase)(nil)

type userSoftDeleter interface {
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type DeleteUserUseCase struct {
	repo userSoftDeleter
}

func NewDeleteUserUseCase(repo userSoftDeleter) *DeleteUserUseCase {
	return &DeleteUserUseCase{repo: repo}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteUser(ctx, id)
}
