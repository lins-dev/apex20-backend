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

type UpdateCampaignInput struct {
	ID          uuid.UUID
	Name        string
	Description *string
}

type CampaignCreator interface {
	Execute(ctx context.Context, input CreateCampaignInput) (campaign.Campaign, error)
}

type CampaignLister interface {
	Execute(ctx context.Context, userID uuid.UUID) ([]campaign.Campaign, error)
}

type CampaignGetter interface {
	Execute(ctx context.Context, id uuid.UUID) (campaign.Campaign, error)
}

type CampaignUpdater interface {
	Execute(ctx context.Context, input UpdateCampaignInput) (campaign.Campaign, error)
}

type CampaignDeleter interface {
	Execute(ctx context.Context, id uuid.UUID) error
}

type MemberInviter interface {
	Execute(ctx context.Context, input InviteMemberInput) (campaign.Member, error)
}

type MemberRemover interface {
	Execute(ctx context.Context, campaignID, userID uuid.UUID) error
}

type MemberRoleGetter interface {
	Execute(ctx context.Context, campaignID, userID uuid.UUID) (permission.Role, error)
}
