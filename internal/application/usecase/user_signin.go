package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/user"
)

var _ port.UserSignIner = (*SignInUseCase)(nil)

type userByEmailGetter interface {
	GetUserByEmail(ctx context.Context, email string) (user.User, error)
}

type SignInUseCase struct {
	repo     userByEmailGetter
	hasher   port.PasswordHasher
	tokenGen port.TokenGenerator
}

func NewSignInUseCase(repo userByEmailGetter, hasher port.PasswordHasher, tokenGen port.TokenGenerator) *SignInUseCase {
	return &SignInUseCase{repo: repo, hasher: hasher, tokenGen: tokenGen}
}

func (uc *SignInUseCase) Execute(ctx context.Context, input port.SignInInput) (port.SignInOutput, error) {
	u, err := uc.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			return port.SignInOutput{}, port.ErrInvalidCredentials
		}
		return port.SignInOutput{}, fmt.Errorf("fetching user: %w", err)
	}

	match, err := uc.hasher.Verify(input.Password, u.PasswordHash)
	if err != nil {
		return port.SignInOutput{}, fmt.Errorf("verifying password: %w", err)
	}
	if !match {
		return port.SignInOutput{}, port.ErrInvalidCredentials
	}

	token, err := uc.tokenGen.Generate(u.ID, u.IsAdmin)
	if err != nil {
		return port.SignInOutput{}, fmt.Errorf("generating token: %w", err)
	}

	return port.SignInOutput{User: u, AccessToken: token}, nil
}
