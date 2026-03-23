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

func (r *PostgresCampaignRepository) CreateCampaignWithMember(ctx context.Context, c campaign.Campaign, m campaign.Member) (err error) {
	tx, txErr := r.db.BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			if err == nil {
				err = rbErr
			}
		}
	}()

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

func (r *PostgresCampaignRepository) ListCampaignsByUserID(ctx context.Context, userID uuid.UUID) ([]campaign.Campaign, error) {
	rows, err := r.queries.ListCampaignsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	campaigns := make([]campaign.Campaign, len(rows))
	for i, row := range rows {
		campaigns[i] = toCampaignDomain(row)
	}
	return campaigns, nil
}

func (r *PostgresCampaignRepository) UpdateCampaign(ctx context.Context, id uuid.UUID, name string, description *string) (campaign.Campaign, error) {
	var desc sql.NullString
	if description != nil {
		desc = sql.NullString{String: *description, Valid: true}
	}
	row, err := r.queries.UpdateCampaign(ctx, repositorygen.UpdateCampaignParams{
		ID:          id,
		Name:        name,
		Description: desc,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return campaign.Campaign{}, port.ErrNotFound
		}
		return campaign.Campaign{}, err
	}
	return toCampaignDomain(row), nil
}

func (r *PostgresCampaignRepository) DeleteCampaign(ctx context.Context, id uuid.UUID) error {
	n, err := r.queries.DeleteCampaign(ctx, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return port.ErrNotFound
	}
	return nil
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
