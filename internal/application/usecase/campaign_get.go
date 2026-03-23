package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/campaign"
)

var _ port.CampaignGetter = (*GetCampaignUseCase)(nil)

type campaignByIDGetter interface {
	GetCampaignByID(ctx context.Context, id uuid.UUID) (campaign.Campaign, error)
}

type GetCampaignUseCase struct {
	repo campaignByIDGetter
}

func NewGetCampaignUseCase(repo campaignByIDGetter) *GetCampaignUseCase {
	return &GetCampaignUseCase{repo: repo}
}

func (uc *GetCampaignUseCase) Execute(ctx context.Context, id uuid.UUID) (campaign.Campaign, error) {
	c, err := uc.repo.GetCampaignByID(ctx, id)
	if err != nil {
		return campaign.Campaign{}, fmt.Errorf("getting campaign: %w", err)
	}
	return c, nil
}
