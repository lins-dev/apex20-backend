package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/campaign"
	"github.com/apex20/backend/internal/domain/permission"
	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository"
	"github.com/apex20/backend/internal/testutil"
)

func createTestUser(t *testing.T, db *sql.DB) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV7()
	require.NoError(t, err)
	_, err = db.ExecContext(context.Background(),
		`INSERT INTO users (id, email, name, password_hash, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		id, id.String()+"@test.com", "Test User", "hash", time.Now(), time.Now(),
	)
	require.NoError(t, err)
	return id
}

func TestPostgresCampaignRepository_CreateWithMember_AndGetByID(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	userID := createTestUser(t, db)
	repo := repository.NewPostgresCampaignRepository(db)
	ctx := context.Background()

	campaignID, err := uuid.NewV7()
	require.NoError(t, err)
	memberID, err := uuid.NewV7()
	require.NoError(t, err)
	now := time.Now()

	c := campaign.Campaign{
		ID:          campaignID,
		UserID:      userID,
		Name:        "Campanha Teste",
		Description: "Descrição",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	m := campaign.Member{
		ID:         memberID,
		CampaignID: campaignID,
		UserID:     userID,
		Role:       permission.RoleGM,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	require.NoError(t, repo.CreateCampaignWithMember(ctx, c, m))

	got, err := repo.GetCampaignByID(ctx, campaignID)
	require.NoError(t, err)
	assert.Equal(t, campaignID, got.ID)
	assert.Equal(t, userID, got.UserID)
	assert.Equal(t, "Campanha Teste", got.Name)
	assert.Equal(t, "Descrição", got.Description)
}

func TestPostgresCampaignRepository_GetByID_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresCampaignRepository(db)

	_, err := repo.GetCampaignByID(context.Background(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestPostgresCampaignMemberRepository_CreateAndGet(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	userID := createTestUser(t, db)
	campaignRepo := repository.NewPostgresCampaignRepository(db)
	memberRepo := repository.NewPostgresCampaignMemberRepository(db)
	ctx := context.Background()

	campaignID, err := uuid.NewV7()
	require.NoError(t, err)
	now := time.Now()

	c := campaign.Campaign{
		ID: campaignID, UserID: userID, Name: "C", CreatedAt: now, UpdatedAt: now,
	}
	gmID, err := uuid.NewV7()
	require.NoError(t, err)
	gm := campaign.Member{
		ID: gmID, CampaignID: campaignID, UserID: userID, Role: permission.RoleGM, CreatedAt: now, UpdatedAt: now,
	}
	require.NoError(t, campaignRepo.CreateCampaignWithMember(ctx, c, gm))

	invitedUserID := createTestUser(t, db)
	memberID, err := uuid.NewV7()
	require.NoError(t, err)
	member := campaign.Member{
		ID:         memberID,
		CampaignID: campaignID,
		UserID:     invitedUserID,
		Role:       permission.RolePlayer,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	require.NoError(t, memberRepo.CreateCampaignMember(ctx, member))

	got, err := memberRepo.GetCampaignMember(ctx, campaignID, invitedUserID)
	require.NoError(t, err)
	assert.Equal(t, memberID, got.ID)
	assert.Equal(t, permission.RolePlayer, got.Role)
}

func TestPostgresCampaignRepository_ListByUserID_ReturnsCampaigns(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	userID := createTestUser(t, db)
	repo := repository.NewPostgresCampaignRepository(db)
	ctx := context.Background()

	for _, name := range []string{"C1", "C2"} {
		cID, _ := uuid.NewV7()
		mID, _ := uuid.NewV7()
		now := time.Now()
		c := campaign.Campaign{ID: cID, UserID: userID, Name: name, CreatedAt: now, UpdatedAt: now}
		m := campaign.Member{ID: mID, CampaignID: cID, UserID: userID, Role: permission.RoleGM, CreatedAt: now, UpdatedAt: now}
		require.NoError(t, repo.CreateCampaignWithMember(ctx, c, m))
	}

	got, err := repo.ListCampaignsByUserID(ctx, userID)

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestPostgresCampaignRepository_ListByUserID_ReturnsEmpty(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresCampaignRepository(db)

	got, err := repo.ListCampaignsByUserID(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestPostgresCampaignRepository_Update_UpdatesFields(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	userID := createTestUser(t, db)
	repo := repository.NewPostgresCampaignRepository(db)
	ctx := context.Background()

	cID, _ := uuid.NewV7()
	mID, _ := uuid.NewV7()
	now := time.Now()
	c := campaign.Campaign{ID: cID, UserID: userID, Name: "Original", Description: "Desc", CreatedAt: now, UpdatedAt: now}
	m := campaign.Member{ID: mID, CampaignID: cID, UserID: userID, Role: permission.RoleGM, CreatedAt: now, UpdatedAt: now}
	require.NoError(t, repo.CreateCampaignWithMember(ctx, c, m))

	updated, err := repo.UpdateCampaign(ctx, cID, "Atualizado", testutil.StrPtr("Nova Desc"))

	require.NoError(t, err)
	assert.Equal(t, "Atualizado", updated.Name)
	assert.Equal(t, "Nova Desc", updated.Description)
	assert.Equal(t, cID, updated.ID)
}

func TestPostgresCampaignRepository_Update_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresCampaignRepository(db)

	_, err := repo.UpdateCampaign(context.Background(), uuid.New(), "X", nil)

	require.Error(t, err)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestPostgresCampaignRepository_Delete_SoftDeletes(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	userID := createTestUser(t, db)
	repo := repository.NewPostgresCampaignRepository(db)
	ctx := context.Background()

	cID, _ := uuid.NewV7()
	mID, _ := uuid.NewV7()
	now := time.Now()
	c := campaign.Campaign{ID: cID, UserID: userID, Name: "C", CreatedAt: now, UpdatedAt: now}
	m := campaign.Member{ID: mID, CampaignID: cID, UserID: userID, Role: permission.RoleGM, CreatedAt: now, UpdatedAt: now}
	require.NoError(t, repo.CreateCampaignWithMember(ctx, c, m))

	require.NoError(t, repo.DeleteCampaign(ctx, cID))

	_, err := repo.GetCampaignByID(ctx, cID)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestPostgresCampaignRepository_Delete_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	repo := repository.NewPostgresCampaignRepository(db)

	err := repo.DeleteCampaign(context.Background(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestPostgresCampaignMemberRepository_GetMember_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanDB(t, db)

	memberRepo := repository.NewPostgresCampaignMemberRepository(db)

	_, err := memberRepo.GetCampaignMember(context.Background(), uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, port.ErrNotFound)
}
