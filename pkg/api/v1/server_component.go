package serverservice

import (
	"time"

	"github.com/google/uuid"

	"go.hollow.sh/serverservice/internal/models"
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

func convertDBServerComponents(dbComponents models.ServerComponentSlice) ([]ServerComponent, error) {
	components := []ServerComponent{}
	if dbComponents == nil {
		return components, nil
	}

	for _, dbC := range dbComponents {
		c := ServerComponent{}
		if err := c.fromDBModel(dbC); err != nil {
			return nil, err
		}

		components = append(components, c)
	}

	return components, nil
}

func (c *ServerComponent) fromDBModel(dbC *models.ServerComponent) error {
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
	c.CreatedAt = dbC.CreatedAt.Time
	c.UpdatedAt = dbC.UpdatedAt.Time

	if dbC.R != nil && dbC.R.ServerComponentType != nil {
		c.ComponentTypeID = dbC.R.ServerComponentType.Slug
		c.ComponentTypeName = dbC.R.ServerComponentType.Name
	}

	attrs, err := convertFromDBAttributes(dbC.R.Attributes)
	if err != nil {
		return err
	}

	c.Attributes = attrs

	return nil
}
