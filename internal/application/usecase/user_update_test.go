package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/user"
	"github.com/apex20/backend/internal/testutil"
)

type stubUserUpdater struct {
	u   user.User
	err error
}

func (s *stubUserUpdater) UpdateUser(_ context.Context, _ uuid.UUID, _ string, _ *string) (user.User, error) {
	return s.u, s.err
}

func TestUpdateUserUseCase_Execute_ReturnsUpdatedUser(t *testing.T) {
	id := uuid.New()
	nick := "dragonslayer"
	expected := user.User{ID: id, Name: "Hero Updated", Nick: nick}
	uc := usecase.NewUpdateUserUseCase(&stubUserUpdater{u: expected})

	u, err := uc.Execute(context.Background(), port.UpdateUserInput{
		ID:   id,
		Name: "Hero Updated",
		Nick: testutil.StrPtr(nick),
	})

	require.NoError(t, err)
	assert.Equal(t, "Hero Updated", u.Name)
	assert.Equal(t, nick, u.Nick)
}

func TestUpdateUserUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewUpdateUserUseCase(&stubUserUpdater{err: port.ErrNotFound})

	_, err := uc.Execute(context.Background(), port.UpdateUserInput{
		ID:   uuid.New(),
		Name: "Hero",
	})

	assert.ErrorIs(t, err, port.ErrNotFound)
}
