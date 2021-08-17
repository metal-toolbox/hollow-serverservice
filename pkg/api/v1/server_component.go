package hollow

import (
	"time"

	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
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

func convertDBServerComponents(dbComponents db.ServerComponentSlice) ([]ServerComponent, error) {
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

func (c *ServerComponent) fromDBModel(dbC *db.ServerComponent) error {
	var err error

	c.UUID, err = uuid.Parse(dbC.ID)
	if err != nil {
		return err
	}

	c.ServerUUID, err = uuid.Parse(dbC.ServerID)
	if err != nil {
		return err
	}

	c.Name = dbC.Name.String
	c.Vendor = dbC.Vendor.String
	c.Model = dbC.Model.String
	c.Serial = dbC.Serial.String
	c.ComponentTypeID = dbC.R.ServerComponentType.Slug
	c.ComponentTypeName = dbC.R.ServerComponentType.Name
	c.CreatedAt = dbC.CreatedAt.Time
	c.UpdatedAt = dbC.UpdatedAt.Time

	attrs, err := convertFromDBAttributes(dbC.R.Attributes)
	if err != nil {
		return err
	}

	c.Attributes = attrs

	return nil
}

// func (c *ServerComponent) toDBModel() (*db.ServerComponent, error) {
// 	dbC := &db.ServerComponent{
// 		ID:       c.UUID.String(),
// 		ServerID: c.ServerUUID.String(),
// 		Name:     null.StringFrom(c.Name),
// 		Vendor:   null.StringFrom(c.Vendor),
// 		Model:    null.StringFrom(c.Model),
// 		Serial:   null.StringFrom(c.Serial),
// 	}

// 	// sct, err := s.FindServerComponentTypeBySlug(c.ComponentTypeID)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// dbC.ServerComponentTypeID = sct.ID

// 	attrs, err := convertToDBAttributes(c.Attributes)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dbC.R.Attributes = attrs

// 	return dbC, nil
// }
