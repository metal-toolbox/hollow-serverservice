package hollow

import (
	"time"

	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// HardwareComponent represents a component of a piece of hardware. These can be
// things like processors, NICs, hard drives, etc.
type HardwareComponent struct {
	UUID                      uuid.UUID    `json:"uuid"`
	HardwareUUID              uuid.UUID    `json:"hardware_uuid"`
	Name                      string       `json:"name"`
	Vendor                    string       `json:"vendor"`
	Model                     string       `json:"model"`
	Serial                    string       `json:"serial"`
	Attributes                []Attributes `json:"attributes"`
	HardwareComponentTypeUUID uuid.UUID    `json:"hardware_component_type_uuid"`
	HardwareComponentTypeName string       `json:"hardware_component_type_name"`
	CreatedAt                 time.Time    `json:"created_at"`
	UpdatedAt                 time.Time    `json:"updated_at"`
}

func convertDBHardwareComponents(dbComponents []db.HardwareComponent) ([]HardwareComponent, error) {
	components := []HardwareComponent{}

	for _, dbC := range dbComponents {
		c := HardwareComponent{}
		if err := c.fromDBModel(dbC); err != nil {
			return nil, err
		}

		components = append(components, c)
	}

	return components, nil
}

func (c *HardwareComponent) fromDBModel(dbC db.HardwareComponent) error {
	c.UUID = dbC.ID
	c.HardwareUUID = dbC.HardwareID
	c.Name = dbC.Name
	c.Vendor = dbC.Vendor
	c.Model = dbC.Model
	c.Serial = dbC.Serial
	c.HardwareComponentTypeUUID = dbC.HardwareComponentType.ID
	c.HardwareComponentTypeName = dbC.HardwareComponentType.Name

	attrs, err := convertDBAttributes(dbC.Attributes)
	if err != nil {
		return err
	}

	c.Attributes = attrs

	return nil
}

func (c *HardwareComponent) toDBModel() (*db.HardwareComponent, error) {
	dbC := &db.HardwareComponent{
		ID:                      c.UUID,
		HardwareID:              c.HardwareUUID,
		Name:                    c.Name,
		Vendor:                  c.Vendor,
		Model:                   c.Model,
		Serial:                  c.Serial,
		HardwareComponentTypeID: c.HardwareComponentTypeUUID,
	}

	// attrs, err := convertDBAttributes(dbC.Attributes)
	// if err != nil {
	// 	return nil, err
	// }

	// c.Attributes = attrs

	// c.HardwareComponentType.fromDBModel(dbC.HardwareComponentType)

	return dbC, nil
}
