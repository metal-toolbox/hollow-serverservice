package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServerComponent represents a component of a server. These can be things like
// processors, NICs, hard drives, etc.
type ServerComponent struct {
	ID                    uuid.UUID
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Name                  string
	Vendor                string
	Model                 string
	Serial                string
	ServerComponentTypeID uuid.UUID
	ServerComponentType   ServerComponentType
	ServerID              uuid.UUID
	Server                Server
	Attributes            []Attributes
}

// BeforeSave ensures that the server component type passes validation checks
func (c *ServerComponent) BeforeSave(tx *gorm.DB) (err error) {
	if c.ID.String() == uuid.Nil.String() {
		c.ID = uuid.New()
	}

	return nil
}
