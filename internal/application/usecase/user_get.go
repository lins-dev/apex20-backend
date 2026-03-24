package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/user"
)

var _ port.UserGetter = (*GetUserUseCase)(nil)

type userByIDGetter interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (user.User, error)
}

type GetUserUseCase struct {
	repo userByIDGetter
}

func NewGetUserUseCase(repo userByIDGetter) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id uuid.UUID) (user.User, error) {
	return uc.repo.GetUserByID(ctx, id)
}
