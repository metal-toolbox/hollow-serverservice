package serverservice

import (
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"go.hollow.sh/serverservice/internal/models"
)

// Firmware represents a firmware file
type Firmware struct {
	UUID        uuid.UUID `json:"uuid"`
	Vendor      string    `json:"vendor"`
	Model       string    `json:"model"`
	Filename    string    `json:"filename"`
	Version     string    `json:"version"`
	ComponentID string    `json:"componentID"`
	Utility     string    `json:"utility"`
	Sha         string    `json:"sha"`
	UpstreamURL string    `json:"upstreamURL"`
}

func (f *Firmware) fromDBModel(dbF *models.Firmware) error {
	var err error

	f.UUID, err = uuid.Parse(dbF.ID)
	if err != nil {
		return err
	}

	f.ComponentID = dbF.ComponentID
	f.Vendor = dbF.Vendor.String
	f.Model = dbF.Model.String
	f.Filename = dbF.Filename.String
	f.Version = dbF.Version.String
	f.Utility = dbF.Utility.String
	f.Sha = dbF.Sha.String
	f.UpstreamURL = dbF.UpstreamURL.String

	return nil
}

func (f *Firmware) toDBModel() (*models.Firmware, error) {
	dbF := &models.Firmware{
		ComponentID: f.ComponentID,
		Vendor:      null.StringFrom(f.Vendor),
		Model:       null.StringFrom(f.Model),
		Filename:    null.StringFrom(f.Filename),
		Version:     null.StringFrom(f.Version),
		Utility:     null.StringFrom(f.Utility),
		Sha:         null.StringFrom(f.Sha),
		UpstreamURL: null.StringFrom(f.UpstreamURL),
	}

	if f.UUID.String() != uuid.Nil.String() {
		dbF.ID = f.UUID.String()
	}

	return dbF, nil
}
