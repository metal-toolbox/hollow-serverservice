package serverservice

import (
	"time"

	"github.com/google/uuid"

	"go.hollow.sh/serverservice/internal/models"
)

// ComponentFirmwareSet represents a group of firmwares
type ComponentFirmwareSet struct {
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
	Name              string                     `json:"name"`
	Attributes        []Attributes               `json:"attributes"`
	ComponentFirmware []ComponentFirmwareVersion `json:"component_firmware"`
	UUID              uuid.UUID                  `json:"uuid"`
}

func (s *ComponentFirmwareSet) fromDBModel(dbFS *models.ComponentFirmwareSet, firmwares []*models.ComponentFirmwareVersion) error {
	var err error

	s.UUID, err = uuid.Parse(dbFS.ID)
	if err != nil {
		return err
	}

	s.Name = dbFS.Name

	for _, firmware := range firmwares {
		f := ComponentFirmwareVersion{}

		err := f.fromDBModel(firmware)
		if err != nil {
			return err
		}

		s.ComponentFirmware = append(s.ComponentFirmware, f)
	}

	// relation attributes
	if dbFS.R.Attributes != nil {
		s.Attributes, err = convertFromDBAttributes(dbFS.R.Attributes)
		if err != nil {
			return err
		}
	}

	s.CreatedAt = dbFS.CreatedAt.Time
	s.UpdatedAt = dbFS.UpdatedAt.Time

	return nil
}

// ComponentFirmwareSetRequest represents the payload to create a firmware set
type ComponentFirmwareSetRequest struct {
	Name                   string       `json:"name"`
	Attributes             []Attributes `json:"attributes"`
	ComponentFirmwareUUIDs []string     `json:"component_firmware_uuids"`
	ID                     uuid.UUID    `json:"uuid"`
}

func (sc *ComponentFirmwareSetRequest) toDBModelFirmwareSet() (*models.ComponentFirmwareSet, error) {
	s := &models.ComponentFirmwareSet{
		ID:   sc.ID.String(),
		Name: sc.Name,
	}

	if sc.ID == uuid.Nil {
		s.ID = uuid.NewString()
	}

	return s, nil
}
