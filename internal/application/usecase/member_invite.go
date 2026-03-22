package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/campaign"
)

var _ port.MemberInviter = (*InviteMemberUseCase)(nil)

type campaignMemberCreator interface {
	CreateCampaignMember(ctx context.Context, m campaign.Member) error
}

type InviteMemberUseCase struct {
	repo campaignMemberCreator
}

func NewInviteMemberUseCase(repo campaignMemberCreator) *InviteMemberUseCase {
	return &InviteMemberUseCase{repo: repo}
}

func (uc *InviteMemberUseCase) Execute(ctx context.Context, input port.InviteMemberInput) (campaign.Member, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return campaign.Member{}, fmt.Errorf("generating member id: %w", err)
	}

	now := time.Now()
	m := campaign.Member{
		ID:         id,
		CampaignID: input.CampaignID,
		UserID:     input.UserID,
		Role:       input.Role,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.repo.CreateCampaignMember(ctx, m); err != nil {
		return campaign.Member{}, fmt.Errorf("inviting member: %w", err)
	}
	return m, nil
}
