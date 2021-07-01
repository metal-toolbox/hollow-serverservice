package db

import (
	"time"

	"github.com/google/uuid"
)

// HardwareComponent represents a component of a piece of hardware. These can be
// things like processors, NICs, hard drives, etc.
type HardwareComponent struct {
	ID                      uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Name                    string
	Vendor                  string
	Model                   string
	Serial                  string
	HardwareComponentTypeID uuid.UUID `gorm:"type:uuid;index"`
	HardwareComponentType   HardwareComponentType
	HardwareID              uuid.UUID `gorm:"type:uuid;index"`
	Hardware                Hardware
	Attributes              []Attributes `gorm:"polymorphic:Entity;"`
}
