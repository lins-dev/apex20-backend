package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/campaign"
	repositorygen "github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository/gen"
)

type PostgresCampaignMemberRepository struct {
	queries *repositorygen.Queries
}

func NewPostgresCampaignMemberRepository(db *sql.DB) *PostgresCampaignMemberRepository {
	return &PostgresCampaignMemberRepository{queries: repositorygen.New(db)}
}

func (r *PostgresCampaignMemberRepository) CreateCampaignMember(ctx context.Context, m campaign.Member) error {
	return r.queries.CreateCampaignMember(ctx, repositorygen.CreateCampaignMemberParams{
		ID:         m.ID,
		CampaignID: m.CampaignID,
		UserID:     m.UserID,
		Role:       m.Role,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	})
}

func (r *PostgresCampaignMemberRepository) GetCampaignMember(ctx context.Context, campaignID, userID uuid.UUID) (campaign.Member, error) {
	row, err := r.queries.GetCampaignMember(ctx, repositorygen.GetCampaignMemberParams{
		CampaignID: campaignID,
		UserID:     userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return campaign.Member{}, port.ErrNotFound
		}
		return campaign.Member{}, err
	}
	return toCampaignMemberDomain(row), nil
}

func (r *PostgresCampaignMemberRepository) DeleteCampaignMember(ctx context.Context, campaignID, userID uuid.UUID) error {
	n, err := r.queries.DeleteCampaignMember(ctx, repositorygen.DeleteCampaignMemberParams{
		CampaignID: campaignID,
		UserID:     userID,
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return port.ErrNotFound
	}
	return nil
}

func toCampaignMemberDomain(row repositorygen.CampaignMember) campaign.Member {
	return campaign.Member{
		ID:         row.ID,
		CampaignID: row.CampaignID,
		UserID:     row.UserID,
		Role:       row.Role,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}
