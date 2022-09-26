package serverservice

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"

	"go.hollow.sh/serverservice/internal/models"
)

// ComponentFirmwareSet represents a group of firmwares
type ComponentFirmwareSet struct {
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
	Name              string                     `json:"name"`
	Metadata          json.RawMessage            `json:"metadata"`
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
	s.Metadata = dbFS.Metadata.JSON

	for _, firmware := range firmwares {
		f := ComponentFirmwareVersion{}

		err := f.fromDBModel(firmware)
		if err != nil {
			return err
		}

		s.ComponentFirmware = append(s.ComponentFirmware, f)
	}

	s.CreatedAt = dbFS.CreatedAt.Time
	s.UpdatedAt = dbFS.UpdatedAt.Time

	return nil
}

// ComponentFirmwareSetRequest represents the payload to create a firmware set
type ComponentFirmwareSetRequest struct {
	Name                   string          `json:"name"`
	Metadata               json.RawMessage `json:"metadata"`
	ComponentFirmwareUUIDs []string        `json:"component_firmware_uuids"`
	ID                     uuid.UUID       `json:"uuid"`
}

func (sc *ComponentFirmwareSetRequest) toDBModelFirmwareSet() (*models.ComponentFirmwareSet, error) {
	s := &models.ComponentFirmwareSet{
		ID:       sc.ID.String(),
		Name:     sc.Name,
		Metadata: null.JSONFrom(sc.Metadata),
	}

	if sc.ID == uuid.Nil {
		s.ID = uuid.NewString()
	}

	return s, nil
}
