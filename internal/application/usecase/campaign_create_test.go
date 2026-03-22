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
	"github.com/apex20/backend/internal/domain/permission"
)

type mockCampaignWithMemberCreator struct {
	err             error
	createdCampaign campaign.Campaign
	createdMember   campaign.Member
}

func (m *mockCampaignWithMemberCreator) CreateCampaignWithMember(_ context.Context, c campaign.Campaign, mem campaign.Member) error {
	if m.err != nil {
		return m.err
	}
	m.createdCampaign = c
	m.createdMember = mem
	return nil
}

func TestCreateCampaignUseCase_Execute_CreatesCampaignWithGMRole(t *testing.T) {
	repo := &mockCampaignWithMemberCreator{}
	uc := usecase.NewCreateCampaignUseCase(repo)
	userID := uuid.New()

	got, err := uc.Execute(context.Background(), port.CreateCampaignInput{
		UserID:      userID,
		Name:        "Campanha Teste",
		Description: "Uma descrição",
	})

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, got.ID)
	assert.Equal(t, userID, got.UserID)
	assert.Equal(t, "Campanha Teste", got.Name)
	assert.Equal(t, "Uma descrição", got.Description)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())

	assert.NotEqual(t, uuid.Nil, repo.createdMember.ID)
	assert.Equal(t, got.ID, repo.createdMember.CampaignID)
	assert.Equal(t, userID, repo.createdMember.UserID)
	assert.Equal(t, permission.RoleGM, repo.createdMember.Role)
	assert.False(t, repo.createdMember.CreatedAt.IsZero())
}

func TestCreateCampaignUseCase_Execute_ReturnsErrorOnRepoFailure(t *testing.T) {
	uc := usecase.NewCreateCampaignUseCase(&mockCampaignWithMemberCreator{err: errors.New("db error")})

	_, err := uc.Execute(context.Background(), port.CreateCampaignInput{
		UserID: uuid.New(),
		Name:   "Test",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "creating campaign")
}
