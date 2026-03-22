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

type mockMemberRoleReader struct {
	result campaign.Member
	err    error
}

func (m *mockMemberRoleReader) GetCampaignMember(_ context.Context, campaignID, userID uuid.UUID) (campaign.Member, error) {
	return m.result, m.err
}

func TestGetMemberRoleUseCase_Execute_ReturnsRole(t *testing.T) {
	campaignID := uuid.New()
	userID := uuid.New()
	repo := &mockMemberRoleReader{
		result: campaign.Member{
			ID:         uuid.New(),
			CampaignID: campaignID,
			UserID:     userID,
			Role:       permission.RoleGM,
		},
	}
	uc := usecase.NewGetMemberRoleUseCase(repo)

	role, err := uc.Execute(context.Background(), campaignID, userID)

	require.NoError(t, err)
	assert.Equal(t, permission.RoleGM, role)
}

func TestGetMemberRoleUseCase_Execute_ReturnsNotFound(t *testing.T) {
	uc := usecase.NewGetMemberRoleUseCase(&mockMemberRoleReader{err: port.ErrNotFound})

	_, err := uc.Execute(context.Background(), uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestGetMemberRoleUseCase_Execute_ReturnsRepoError(t *testing.T) {
	uc := usecase.NewGetMemberRoleUseCase(&mockMemberRoleReader{err: errors.New("db error")})

	_, err := uc.Execute(context.Background(), uuid.New(), uuid.New())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting member role")
}
