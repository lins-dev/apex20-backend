package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/campaign"
)

type mockCampaignsByUserIDLister struct {
	campaigns []campaign.Campaign
	err       error
}

func (m *mockCampaignsByUserIDLister) ListCampaignsByUserID(_ context.Context, _ uuid.UUID) ([]campaign.Campaign, error) {
	return m.campaigns, m.err
}

func TestListCampaignsUseCase_Execute_ReturnsCampaigns(t *testing.T) {
	userID := uuid.New()
	expected := []campaign.Campaign{
		{ID: uuid.New(), UserID: userID, Name: "C1"},
		{ID: uuid.New(), UserID: userID, Name: "C2"},
	}
	uc := usecase.NewListCampaignsUseCase(&mockCampaignsByUserIDLister{campaigns: expected})

	got, err := uc.Execute(context.Background(), userID)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestListCampaignsUseCase_Execute_ReturnsEmptySlice(t *testing.T) {
	uc := usecase.NewListCampaignsUseCase(&mockCampaignsByUserIDLister{campaigns: nil})

	got, err := uc.Execute(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestListCampaignsUseCase_Execute_ReturnsErrorOnRepoFailure(t *testing.T) {
	uc := usecase.NewListCampaignsUseCase(&mockCampaignsByUserIDLister{err: errors.New("db error")})

	_, err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "listing campaigns")
}
