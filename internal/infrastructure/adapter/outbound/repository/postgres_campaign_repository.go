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

type PostgresCampaignRepository struct {
	db      *sql.DB
	queries *repositorygen.Queries
}

func NewPostgresCampaignRepository(db *sql.DB) *PostgresCampaignRepository {
	return &PostgresCampaignRepository{db: db, queries: repositorygen.New(db)}
}

func (r *PostgresCampaignRepository) CreateCampaignWithMember(ctx context.Context, c campaign.Campaign, m campaign.Member) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := r.queries.WithTx(tx)

	if err := q.CreateCampaign(ctx, repositorygen.CreateCampaignParams{
		ID:          c.ID,
		UserID:      c.UserID,
		Name:        c.Name,
		Description: sql.NullString{String: c.Description, Valid: c.Description != ""},
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}); err != nil {
		return err
	}

	if err := q.CreateCampaignMember(ctx, repositorygen.CreateCampaignMemberParams{
		ID:         m.ID,
		CampaignID: m.CampaignID,
		UserID:     m.UserID,
		Role:       m.Role,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresCampaignRepository) GetCampaignByID(ctx context.Context, id uuid.UUID) (campaign.Campaign, error) {
	row, err := r.queries.GetCampaignByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return campaign.Campaign{}, port.ErrNotFound
		}
		return campaign.Campaign{}, err
	}
	return toCampaignDomain(row), nil
}

func toCampaignDomain(row repositorygen.Campaign) campaign.Campaign {
	return campaign.Campaign{
		ID:          row.ID,
		UserID:      row.UserID,
		Name:        row.Name,
		Description: row.Description.String,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
