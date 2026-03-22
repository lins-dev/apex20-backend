package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/campaign"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.CampaignCreator = (*CreateCampaignUseCase)(nil)

type campaignWithMemberCreator interface {
	CreateCampaignWithMember(ctx context.Context, c campaign.Campaign, m campaign.Member) error
}

type CreateCampaignUseCase struct {
	repo campaignWithMemberCreator
}

func NewCreateCampaignUseCase(repo campaignWithMemberCreator) *CreateCampaignUseCase {
	return &CreateCampaignUseCase{repo: repo}
}

func (uc *CreateCampaignUseCase) Execute(ctx context.Context, input port.CreateCampaignInput) (campaign.Campaign, error) {
	campaignID, err := uuid.NewV7()
	if err != nil {
		return campaign.Campaign{}, fmt.Errorf("generating campaign id: %w", err)
	}
	memberID, err := uuid.NewV7()
	if err != nil {
		return campaign.Campaign{}, fmt.Errorf("generating member id: %w", err)
	}

	now := time.Now()
	c := campaign.Campaign{
		ID:          campaignID,
		UserID:      input.UserID,
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	m := campaign.Member{
		ID:         memberID,
		CampaignID: campaignID,
		UserID:     input.UserID,
		Role:       permission.RoleGM,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.repo.CreateCampaignWithMember(ctx, c, m); err != nil {
		return campaign.Campaign{}, fmt.Errorf("creating campaign: %w", err)
	}
	return c, nil
}
