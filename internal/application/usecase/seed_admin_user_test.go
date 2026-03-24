package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/user"
)

func TestSeedAdminUserUseCase_Execute_SkipsWhenAdminAlreadyExists(t *testing.T) {
	getter := &stubUserByEmailGetter{u: user.User{Email: "admin@apex20.com"}}
	creator := &trackingUserCreator{}
	hasher := &stubPasswordHasher{hash: "hashed"}

	uc := usecase.NewSeedAdminUserUseCase(getter, creator, hasher)
	err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.False(t, creator.called)
}

func TestSeedAdminUserUseCase_Execute_CreatesAdminWhenNotExists(t *testing.T) {
	getter := &stubUserByEmailGetter{err: port.ErrNotFound}
	creator := &trackingUserCreator{}
	hasher := &stubPasswordHasher{hash: "hashed_pw"}

	uc := usecase.NewSeedAdminUserUseCase(getter, creator, hasher)
	err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.True(t, creator.called)
	assert.Equal(t, "admin@apex20.com", creator.created.Email)
	assert.True(t, creator.created.IsAdmin)
	assert.Equal(t, "hashed_pw", creator.created.PasswordHash)
}

// stubUserCreator already defined in user_signup_test.go — extend with tracking
type trackingUserCreator struct {
	called  bool
	created user.User
	err     error
}

func (s *trackingUserCreator) CreateUser(_ context.Context, u user.User) error {
	s.called = true
	s.created = u
	return s.err
}
