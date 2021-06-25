package hollow

import (
	"time"

	"github.com/google/uuid"
)

// Hardware provides a versioned representation of hardware from the datastore
type Hardware struct {
	UUID         uuid.UUID `json:"uuid"`
	FacilityCode string    `json:"facility"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
