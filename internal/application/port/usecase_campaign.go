package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/domain/campaign"
	"github.com/apex20/backend/internal/domain/permission"
)

type CreateCampaignInput struct {
	UserID      uuid.UUID
	Name        string
	Description string
}

type InviteMemberInput struct {
	CampaignID uuid.UUID
	UserID     uuid.UUID
	Role       permission.Role
}

type CampaignCreator interface {
	Execute(ctx context.Context, input CreateCampaignInput) (campaign.Campaign, error)
}

type MemberInviter interface {
	Execute(ctx context.Context, input InviteMemberInput) (campaign.Member, error)
}

type MemberRoleGetter interface {
	Execute(ctx context.Context, campaignID, userID uuid.UUID) (permission.Role, error)
}
