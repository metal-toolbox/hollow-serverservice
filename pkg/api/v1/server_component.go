package hollow

import (
	"time"

	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/gormdb"
)

// ServerComponent represents a component of a server. These can be things like
// processors, NICs, hard drives, etc.
type ServerComponent struct {
	UUID              uuid.UUID    `json:"uuid"`
	ServerUUID        uuid.UUID    `json:"server_uuid"`
	Name              string       `json:"name"`
	Vendor            string       `json:"vendor"`
	Model             string       `json:"model"`
	Serial            string       `json:"serial"`
	Attributes        []Attributes `json:"attributes"`
	ComponentTypeID   string       `json:"component_type_id"`
	ComponentTypeName string       `json:"component_type_name"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

func convertDBServerComponents(dbComponents []gormdb.ServerComponent) ([]ServerComponent, error) {
	components := []ServerComponent{}

	for _, dbC := range dbComponents {
		c := ServerComponent{}
		if err := c.fromDBModel(dbC); err != nil {
			return nil, err
		}

		components = append(components, c)
	}

	return components, nil
}

func (c *ServerComponent) fromDBModel(dbC gormdb.ServerComponent) error {
	c.UUID = dbC.ID
	c.ServerUUID = dbC.ServerID
	c.Name = dbC.Name
	c.Vendor = dbC.Vendor
	c.Model = dbC.Model
	c.Serial = dbC.Serial
	c.ComponentTypeID = dbC.ServerComponentType.Slug
	c.ComponentTypeName = dbC.ServerComponentType.Name
	c.CreatedAt = dbC.CreatedAt
	c.UpdatedAt = dbC.UpdatedAt

	attrs, err := convertFromDBAttributes(dbC.Attributes)
	if err != nil {
		return err
	}

	c.Attributes = attrs

	return nil
}

func (c *ServerComponent) toDBModel(s *gormdb.Store) (*gormdb.ServerComponent, error) {
	dbC := &gormdb.ServerComponent{
		ID:       c.UUID,
		ServerID: c.ServerUUID,
		Name:     c.Name,
		Vendor:   c.Vendor,
		Model:    c.Model,
		Serial:   c.Serial,
	}

	sct, err := s.FindServerComponentTypeBySlug(c.ComponentTypeID)
	if err != nil {
		return nil, err
	}

	dbC.ServerComponentTypeID = sct.ID

	attrs, err := convertToDBAttributes(c.Attributes)
	if err != nil {
		return nil, err
	}

	dbC.Attributes = attrs

	return dbC, nil
}
