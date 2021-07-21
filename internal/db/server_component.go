package db

import (
	"time"

	"github.com/google/uuid"
)

// ServerComponent represents a component of a server. These can be things like
// processors, NICs, hard drives, etc.
type ServerComponent struct {
	ID                    uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Name                  string
	Vendor                string
	Model                 string
	Serial                string
	ServerComponentTypeID uuid.UUID `gorm:"type:uuid;index"`
	ServerComponentType   ServerComponentType
	ServerID              uuid.UUID `gorm:"type:uuid;index"`
	Server                Server
	Attributes            []Attributes `gorm:"polymorphic:Entity;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
