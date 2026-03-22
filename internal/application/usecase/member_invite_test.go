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

type mockMemberCreator struct {
	err     error
	created []campaign.Member
}

func (m *mockMemberCreator) CreateCampaignMember(_ context.Context, mem campaign.Member) error {
	if m.err != nil {
		return m.err
	}
	m.created = append(m.created, mem)
	return nil
}

func TestInviteMemberUseCase_Execute_CreatesPlayerMember(t *testing.T) {
	repo := &mockMemberCreator{}
	uc := usecase.NewInviteMemberUseCase(repo)
	campaignID := uuid.New()
	userID := uuid.New()

	got, err := uc.Execute(context.Background(), port.InviteMemberInput{
		CampaignID: campaignID,
		UserID:     userID,
		Role:       permission.RolePlayer,
	})

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, got.ID)
	assert.Equal(t, campaignID, got.CampaignID)
	assert.Equal(t, userID, got.UserID)
	assert.Equal(t, permission.RolePlayer, got.Role)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())
	require.Len(t, repo.created, 1)
	assert.Equal(t, got.ID, repo.created[0].ID)
}

func TestInviteMemberUseCase_Execute_CreatesTrustedMember(t *testing.T) {
	repo := &mockMemberCreator{}
	uc := usecase.NewInviteMemberUseCase(repo)

	got, err := uc.Execute(context.Background(), port.InviteMemberInput{
		CampaignID: uuid.New(),
		UserID:     uuid.New(),
		Role:       permission.RoleTrusted,
	})

	require.NoError(t, err)
	assert.Equal(t, permission.RoleTrusted, got.Role)
}

func TestInviteMemberUseCase_Execute_ReturnsErrorOnRepoFailure(t *testing.T) {
	uc := usecase.NewInviteMemberUseCase(&mockMemberCreator{err: errors.New("duplicate key")})

	_, err := uc.Execute(context.Background(), port.InviteMemberInput{
		CampaignID: uuid.New(),
		UserID:     uuid.New(),
		Role:       permission.RolePlayer,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "inviting member")
}
