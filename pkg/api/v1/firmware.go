package serverservice

import (
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"

	"go.hollow.sh/serverservice/internal/models"
)

// ComponentFirmwareVersion represents a firmware file
type ComponentFirmwareVersion struct {
	UUID        uuid.UUID `json:"uuid"`
	Vendor      string    `json:"vendor"`
	Model       string    `json:"model"`
	Filename    string    `json:"filename"`
	Version     string    `json:"version"`
	Component   string    `json:"component"`
	Utility     string    `json:"utility"`
	Sha         string    `json:"sha"`
	UpstreamURL string    `json:"upstreamURL"`
}

func (f *ComponentFirmwareVersion) fromDBModel(dbF *models.ComponentFirmwareVersion) error {
	var err error

	f.UUID, err = uuid.Parse(dbF.ID)
	if err != nil {
		return err
	}

	f.Component = dbF.Component.String
	f.Vendor = dbF.Vendor.String
	f.Model = dbF.Model.String
	f.Filename = dbF.Filename.String
	f.Version = dbF.Version.String
	f.Utility = dbF.Utility.String
	f.Sha = dbF.Sha.String
	f.UpstreamURL = dbF.UpstreamURL.String

	return nil
}

func (f *ComponentFirmwareVersion) toDBModel() (*models.ComponentFirmwareVersion, error) {
	dbF := &models.ComponentFirmwareVersion{
		Component:   null.StringFrom(f.Component),
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
