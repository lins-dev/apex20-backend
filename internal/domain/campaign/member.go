package campaign

import (
	"time"

	"github.com/google/uuid"

	"github.com/apex20/backend/internal/domain/permission"
)

// Member represents a user's membership and role within a specific campaign.
type Member struct {
	ID         uuid.UUID
	CampaignID uuid.UUID
	UserID     uuid.UUID
	Role       permission.Role
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
