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
	"github.com/apex20/backend/internal/domain/campaign"
	"github.com/apex20/backend/internal/testutil"
)

type mockCampaignUpdater struct {
	updated campaign.Campaign
	err     error
}

func (m *mockCampaignUpdater) UpdateCampaign(_ context.Context, id uuid.UUID, name string, description *string) (campaign.Campaign, error) {
	if m.err != nil {
		return campaign.Campaign{}, m.err
	}
	desc := ""
	if description != nil {
		desc = *description
	}
	m.updated = campaign.Campaign{ID: id, Name: name, Description: desc}
	return m.updated, nil
}

func TestUpdateCampaignUseCase_Execute_ReturnsUpdatedCampaign(t *testing.T) {
	id := uuid.New()
	uc := usecase.NewUpdateCampaignUseCase(&mockCampaignUpdater{})

	got, err := uc.Execute(context.Background(), port.UpdateCampaignInput{
		ID:          id,
		Name:        "Novo Nome",
		Description: testutil.StrPtr("Nova Desc"),
	})

	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "Novo Nome", got.Name)
	assert.Equal(t, "Nova Desc", got.Description)
}

func TestUpdateCampaignUseCase_Execute_ReturnsUpdatedCampaign_NilDescription(t *testing.T) {
	id := uuid.New()
	uc := usecase.NewUpdateCampaignUseCase(&mockCampaignUpdater{})

	got, err := uc.Execute(context.Background(), port.UpdateCampaignInput{
		ID:          id,
		Name:        "Novo Nome",
		Description: nil,
	})

	require.NoError(t, err)
	assert.Equal(t, "", got.Description)
}

func TestUpdateCampaignUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewUpdateCampaignUseCase(&mockCampaignUpdater{err: port.ErrNotFound})

	_, err := uc.Execute(context.Background(), port.UpdateCampaignInput{ID: uuid.New(), Name: "X"})

	require.Error(t, err)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestUpdateCampaignUseCase_Execute_ReturnsErrorOnRepoFailure(t *testing.T) {
	uc := usecase.NewUpdateCampaignUseCase(&mockCampaignUpdater{err: errors.New("db error")})

	_, err := uc.Execute(context.Background(), port.UpdateCampaignInput{ID: uuid.New(), Name: "X"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "updating campaign")
}
