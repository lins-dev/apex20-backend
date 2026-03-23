package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
)

var _ port.CampaignDeleter = (*DeleteCampaignUseCase)(nil)

type campaignDeleterRepo interface {
	DeleteCampaign(ctx context.Context, id uuid.UUID) error
}

type DeleteCampaignUseCase struct {
	repo campaignDeleterRepo
}

func NewDeleteCampaignUseCase(repo campaignDeleterRepo) *DeleteCampaignUseCase {
	return &DeleteCampaignUseCase{repo: repo}
}

func (uc *DeleteCampaignUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.DeleteCampaign(ctx, id); err != nil {
		return fmt.Errorf("deleting campaign: %w", err)
	}
	return nil
}
