package hollow

import (
	"time"

	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// Hardware represents a piece of hardware in a facility. These are the
// details of the physical hardware
type Hardware struct {
	UUID                uuid.UUID             `json:"uuid"`
	FacilityCode        string                `json:"facility"`
	Attributes          []Attributes          `json:"attributes"`
	HardwareComponents  []HardwareComponent   `json:"hardware_components"`
	VersionedAttributes []VersionedAttributes `json:"versioned_attributes"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}

func (h *Hardware) fromDBModel(dbH db.Hardware) error {
	var err error

	h.UUID = dbH.ID
	h.FacilityCode = dbH.FacilityCode
	h.CreatedAt = dbH.CreatedAt
	h.UpdatedAt = dbH.UpdatedAt

	h.Attributes, err = convertFromDBAttributes(dbH.Attributes)
	if err != nil {
		return err
	}

	h.HardwareComponents, err = convertDBHardwareComponents(dbH.HardwareComponents)
	if err != nil {
		return err
	}

	h.VersionedAttributes, err = convertFromDBVersionedAttributes(dbH.VersionedAttributes)
	if err != nil {
		return err
	}

	return nil
}

func (h *Hardware) toDBModel() (*db.Hardware, error) {
	dbC := &db.Hardware{
		ID:           h.UUID,
		FacilityCode: h.FacilityCode,
	}

	for _, hc := range h.HardwareComponents {
		c, err := hc.toDBModel()
		if err != nil {
			return nil, err
		}

		dbC.HardwareComponents = append(dbC.HardwareComponents, *c)
	}

	attrs, err := convertToDBAttributes(h.Attributes)
	if err != nil {
		return nil, err
	}

	dbC.Attributes = attrs

	verAttrs, err := convertToDBVersionedAttributes(h.VersionedAttributes)
	if err != nil {
		return nil, err
	}

	dbC.VersionedAttributes = verAttrs

	return dbC, nil
}
