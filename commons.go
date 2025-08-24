package genesis

import (
	"time"

	"github.com/pborman/uuid"
)

// NextUUID generates a new UUID.
func NextUUID() string {
	return uuid.New()
}

// AuditFields holds common fields for tracking record changes.
type AuditFields struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
