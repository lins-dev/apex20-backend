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

type mockCampaignDeleter struct {
	err error
}

func (m *mockCampaignDeleter) DeleteCampaign(_ context.Context, _ uuid.UUID) error {
	return m.err
}

func TestDeleteCampaignUseCase_Execute_Succeeds(t *testing.T) {
	uc := usecase.NewDeleteCampaignUseCase(&mockCampaignDeleter{})

	err := uc.Execute(context.Background(), uuid.New())

	require.NoError(t, err)
}

func TestDeleteCampaignUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewDeleteCampaignUseCase(&mockCampaignDeleter{err: port.ErrNotFound})

	err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestDeleteCampaignUseCase_Execute_ReturnsErrorOnRepoFailure(t *testing.T) {
	uc := usecase.NewDeleteCampaignUseCase(&mockCampaignDeleter{err: errors.New("db error")})

	err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "deleting campaign")
}
