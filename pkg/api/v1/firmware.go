package serverservice

import (
	"github.com/google/uuid"

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
	Checksum    string    `json:"checksum"`
	UpstreamURL string    `json:"upstreamURL"`
	S3URL       string    `json:"s3URL"`
}

func (f *ComponentFirmwareVersion) fromDBModel(dbF *models.ComponentFirmwareVersion) error {
	var err error

	f.UUID, err = uuid.Parse(dbF.ID)
	if err != nil {
		return err
	}

	f.Component = dbF.Component
	f.Vendor = dbF.Vendor
	f.Model = dbF.Model
	f.Filename = dbF.Filename
	f.Version = dbF.Version
	f.Utility = dbF.Utility
	f.Checksum = dbF.Checksum
	f.UpstreamURL = dbF.UpstreamURL
	f.S3URL = dbF.S3URL

	return nil
}

func (f *ComponentFirmwareVersion) toDBModel() (*models.ComponentFirmwareVersion, error) {
	dbF := &models.ComponentFirmwareVersion{
		Component:   f.Component,
		Vendor:      f.Vendor,
		Model:       f.Model,
		Filename:    f.Filename,
		Version:     f.Version,
		Utility:     f.Utility,
		Checksum:    f.Checksum,
		UpstreamURL: f.UpstreamURL,
		S3URL:       f.S3URL,
	}

	if f.UUID.String() != uuid.Nil.String() {
		dbF.ID = f.UUID.String()
	}

	return dbF, nil
}
