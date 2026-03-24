package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/user"
)

const (
	adminEmail    = "admin@apex20.com"
	adminName     = "Admin"
	adminPassword = "admin123"
)

type SeedAdminUserUseCase struct {
	getter  userByEmailGetter
	creator userCreator
	hasher  port.PasswordHasher
}

func NewSeedAdminUserUseCase(getter userByEmailGetter, creator userCreator, hasher port.PasswordHasher) *SeedAdminUserUseCase {
	return &SeedAdminUserUseCase{getter: getter, creator: creator, hasher: hasher}
}

// Execute creates the default admin user if it does not already exist. Idempotent.
func (uc *SeedAdminUserUseCase) Execute(ctx context.Context) error {
	_, err := uc.getter.GetUserByEmail(ctx, adminEmail)
	if err == nil {
		return nil
	}
	if !errors.Is(err, port.ErrNotFound) {
		return fmt.Errorf("checking admin user: %w", err)
	}

	hash, err := uc.hasher.Hash(adminPassword)
	if err != nil {
		return fmt.Errorf("hashing admin password: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generating admin id: %w", err)
	}

	now := time.Now()
	admin := user.User{
		ID:           id,
		Email:        adminEmail,
		Name:         adminName,
		PasswordHash: hash,
		IsAdmin:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.creator.CreateUser(ctx, admin); err != nil {
		return fmt.Errorf("creating admin user: %w", err)
	}

	return nil
}
