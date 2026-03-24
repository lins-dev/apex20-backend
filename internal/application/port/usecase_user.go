package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/domain/user"
)

type SignUpInput struct {
	Email    string
	Password string
	Name     string
}

type SignUpOutput struct {
	User        user.User
	AccessToken string
}

type SignInInput struct {
	Email    string
	Password string
}

type SignInOutput struct {
	User        user.User
	AccessToken string
}

type UpdateUserInput struct {
	ID   uuid.UUID
	Name string
	Nick *string
}

type UserSignUpper interface {
	Execute(ctx context.Context, input SignUpInput) (SignUpOutput, error)
}

type UserSignIner interface {
	Execute(ctx context.Context, input SignInInput) (SignInOutput, error)
}

type UserGetter interface {
	Execute(ctx context.Context, id uuid.UUID) (user.User, error)
}

type UserUpdater interface {
	Execute(ctx context.Context, input UpdateUserInput) (user.User, error)
}

type UserDeleter interface {
	Execute(ctx context.Context, id uuid.UUID) error
}
