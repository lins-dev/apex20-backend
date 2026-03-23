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
)

type mockCampaignByIDGetter struct {
	campaign campaign.Campaign
	err      error
}

func (m *mockCampaignByIDGetter) GetCampaignByID(_ context.Context, _ uuid.UUID) (campaign.Campaign, error) {
	return m.campaign, m.err
}

func TestGetCampaignUseCase_Execute_ReturnsCampaign(t *testing.T) {
	id := uuid.New()
	expected := campaign.Campaign{ID: id, Name: "Test"}
	uc := usecase.NewGetCampaignUseCase(&mockCampaignByIDGetter{campaign: expected})

	got, err := uc.Execute(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetCampaignUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewGetCampaignUseCase(&mockCampaignByIDGetter{err: port.ErrNotFound})

	_, err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestGetCampaignUseCase_Execute_ReturnsErrorOnRepoFailure(t *testing.T) {
	uc := usecase.NewGetCampaignUseCase(&mockCampaignByIDGetter{err: errors.New("db error")})

	_, err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting campaign")
}
