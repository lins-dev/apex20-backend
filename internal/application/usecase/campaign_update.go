package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/campaign"
)

var _ port.CampaignUpdater = (*UpdateCampaignUseCase)(nil)

type campaignUpdaterRepo interface {
	UpdateCampaign(ctx context.Context, id uuid.UUID, name string, description *string) (campaign.Campaign, error)
}

type UpdateCampaignUseCase struct {
	repo campaignUpdaterRepo
}

func NewUpdateCampaignUseCase(repo campaignUpdaterRepo) *UpdateCampaignUseCase {
	return &UpdateCampaignUseCase{repo: repo}
}

func (uc *UpdateCampaignUseCase) Execute(ctx context.Context, input port.UpdateCampaignInput) (campaign.Campaign, error) {
	c, err := uc.repo.UpdateCampaign(ctx, input.ID, input.Name, input.Description)
	if err != nil {
		return campaign.Campaign{}, fmt.Errorf("updating campaign: %w", err)
	}
	return c, nil
}
