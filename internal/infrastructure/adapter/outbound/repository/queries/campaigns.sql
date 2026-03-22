-- name: CreateCampaign :exec
INSERT INTO campaigns (id, user_id, name, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetCampaignByID :one
SELECT id, user_id, name, description, created_at, updated_at, deleted_at
FROM campaigns
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCampaignsByUserID :many
SELECT id, user_id, name, description, created_at, updated_at, deleted_at
FROM campaigns
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;
