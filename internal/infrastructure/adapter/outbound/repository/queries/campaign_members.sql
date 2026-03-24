-- name: GetCampaignMember :one
SELECT * FROM campaign_members
WHERE campaign_id = $1 AND user_id = $2;

-- name: ListCampaignMembers :many
SELECT * FROM campaign_members
WHERE campaign_id = $1
ORDER BY created_at;

-- name: CreateCampaignMember :exec
INSERT INTO campaign_members (id, campaign_id, user_id, role, display_name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateCampaignMemberDisplayName :exec
UPDATE campaign_members
SET display_name = $3, updated_at = NOW()
WHERE campaign_id = $1 AND user_id = $2;

-- name: UpdateCampaignMemberRole :exec
UPDATE campaign_members SET role = $3, updated_at = $4
WHERE campaign_id = $1 AND user_id = $2;

-- name: DeleteCampaignMember :execrows
DELETE FROM campaign_members
WHERE campaign_id = $1 AND user_id = $2;
