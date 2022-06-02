package serverservice

import (
	"time"

	"github.com/google/uuid"

	"go.hollow.sh/serverservice/internal/models"
)

// ServerComponent represents a component of a server. These can be things like
// processors, NICs, hard drives, etc.
//
// Note: when setting validator struct tags, ensure no extra spaces are present between
//       comma separated values or validation will fail with a not so useful 500 error.
type ServerComponent struct {
	UUID                uuid.UUID             `json:"uuid"`
	ServerUUID          uuid.UUID             `json:"server_uuid" binding:"required"`
	Name                string                `json:"name" binding:"required"`
	Vendor              string                `json:"vendor" binding:"required,lowercase"`
	Model               string                `json:"model" binding:"required,lowercase"`
	Serial              string                `json:"serial" binding:"required,lowercase"`
	Attributes          []Attributes          `json:"attributes"`
	VersionedAttributes []VersionedAttributes `json:"versioned_attributes"`
	ComponentTypeID     string                `json:"component_type_id" binding:"required"`
	ComponentTypeName   string                `json:"component_type_name" binding:"required"`
	ComponentTypeSlug   string                `json:"component_type_slug"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}

// ServerComponentSlice is a slice of ServerComponent objects
type ServerComponentSlice []ServerComponent

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
		c.ComponentTypeID = dbC.R.ServerComponentType.ID
		c.ComponentTypeName = dbC.R.ServerComponentType.Name
		c.ComponentTypeSlug = dbC.R.ServerComponentType.Slug
	}

	attrs, err := convertFromDBAttributes(dbC.R.Attributes)
	if err != nil {
		return err
	}

	c.Attributes = attrs

	return nil
}
