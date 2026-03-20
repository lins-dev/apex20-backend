package campaign

import (
	"time"

	apexv1 "github.com/apex20/contracts/proto/apex20/v1"
	"github.com/google/uuid"
)

// Role is the campaign-scoped role type from Protobuf.
// Valid values: ROLE_GM, ROLE_PLAYER, ROLE_TRUSTED.
type Role = apexv1.Role

const (
	RoleGM      = apexv1.Role_ROLE_GM
	RolePlayer  = apexv1.Role_ROLE_PLAYER
	RoleTrusted = apexv1.Role_ROLE_TRUSTED
)

// Member represents a user's membership and role within a specific campaign.
type Member struct {
	ID         uuid.UUID
	CampaignID uuid.UUID
	UserID     uuid.UUID
	Role       Role
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
