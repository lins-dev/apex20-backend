package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/user"
)

var _ port.UserSignUpper = (*SignUpUseCase)(nil)

type userCreator interface {
	CreateUser(ctx context.Context, u user.User) error
}

type SignUpUseCase struct {
	repo      userCreator
	hasher    port.PasswordHasher
	tokenGen  port.TokenGenerator
}

func NewSignUpUseCase(repo userCreator, hasher port.PasswordHasher, tokenGen port.TokenGenerator) *SignUpUseCase {
	return &SignUpUseCase{repo: repo, hasher: hasher, tokenGen: tokenGen}
}

func (uc *SignUpUseCase) Execute(ctx context.Context, input port.SignUpInput) (port.SignUpOutput, error) {
	hash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return port.SignUpOutput{}, fmt.Errorf("hashing password: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return port.SignUpOutput{}, fmt.Errorf("generating user id: %w", err)
	}

	now := time.Now()
	u := user.User{
		ID:           id,
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: hash,
		IsAdmin:      false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.repo.CreateUser(ctx, u); err != nil {
		return port.SignUpOutput{}, err
	}

	token, err := uc.tokenGen.Generate(u.ID, u.IsAdmin)
	if err != nil {
		return port.SignUpOutput{}, fmt.Errorf("generating token: %w", err)
	}

	return port.SignUpOutput{User: u, AccessToken: token}, nil
}
