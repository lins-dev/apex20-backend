package campaign

import (
	"time"

	"github.com/google/uuid"
)

type Campaign struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
