package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
)

type mockCampaignMemberDeleter struct {
	err error
}

func (m *mockCampaignMemberDeleter) DeleteCampaignMember(_ context.Context, _, _ uuid.UUID) error {
	return m.err
}

func TestRemoveMemberUseCase_Execute_Succeeds(t *testing.T) {
	uc := usecase.NewRemoveMemberUseCase(&mockCampaignMemberDeleter{})

	err := uc.Execute(context.Background(), uuid.New(), uuid.New())

	require.NoError(t, err)
}

func TestRemoveMemberUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewRemoveMemberUseCase(&mockCampaignMemberDeleter{err: port.ErrNotFound})

	err := uc.Execute(context.Background(), uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestRemoveMemberUseCase_Execute_ReturnsErrorOnRepoFailure(t *testing.T) {
	uc := usecase.NewRemoveMemberUseCase(&mockCampaignMemberDeleter{err: errors.New("db error")})

	err := uc.Execute(context.Background(), uuid.New(), uuid.New())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "removing member")
}
