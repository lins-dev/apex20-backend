package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/campaign"
	"github.com/apex20/backend/internal/domain/permission"
)

var _ port.MemberRoleGetter = (*GetMemberRoleUseCase)(nil)

type campaignMemberReader interface {
	GetCampaignMember(ctx context.Context, campaignID, userID uuid.UUID) (campaign.Member, error)
}

type GetMemberRoleUseCase struct {
	repo campaignMemberReader
}

func NewGetMemberRoleUseCase(repo campaignMemberReader) *GetMemberRoleUseCase {
	return &GetMemberRoleUseCase{repo: repo}
}

func (uc *GetMemberRoleUseCase) Execute(ctx context.Context, campaignID, userID uuid.UUID) (permission.Role, error) {
	m, err := uc.repo.GetCampaignMember(ctx, campaignID, userID)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			return 0, port.ErrNotFound
		}
		return 0, fmt.Errorf("getting member role: %w", err)
	}
	return m.Role, nil
}
