package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/campaign"
)

var _ port.CampaignLister = (*ListCampaignsUseCase)(nil)

type campaignsByUserIDLister interface {
	ListCampaignsByUserID(ctx context.Context, userID uuid.UUID) ([]campaign.Campaign, error)
}

type ListCampaignsUseCase struct {
	repo campaignsByUserIDLister
}

func NewListCampaignsUseCase(repo campaignsByUserIDLister) *ListCampaignsUseCase {
	return &ListCampaignsUseCase{repo: repo}
}

func (uc *ListCampaignsUseCase) Execute(ctx context.Context, userID uuid.UUID) ([]campaign.Campaign, error) {
	campaigns, err := uc.repo.ListCampaignsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listing campaigns: %w", err)
	}
	return campaigns, nil
}
