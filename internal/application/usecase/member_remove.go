package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
)

var _ port.MemberRemover = (*RemoveMemberUseCase)(nil)

type campaignMemberDeleter interface {
	DeleteCampaignMember(ctx context.Context, campaignID, userID uuid.UUID) error
}

type RemoveMemberUseCase struct {
	repo campaignMemberDeleter
}

func NewRemoveMemberUseCase(repo campaignMemberDeleter) *RemoveMemberUseCase {
	return &RemoveMemberUseCase{repo: repo}
}

func (uc *RemoveMemberUseCase) Execute(ctx context.Context, campaignID, userID uuid.UUID) error {
	if err := uc.repo.DeleteCampaignMember(ctx, campaignID, userID); err != nil {
		return fmt.Errorf("removing member: %w", err)
	}
	return nil
}
