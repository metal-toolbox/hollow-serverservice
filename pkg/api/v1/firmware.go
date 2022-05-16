package serverservice

import (
	"github.com/google/uuid"

	"go.hollow.sh/serverservice/internal/models"
)

// ComponentFirmwareVersion represents a firmware file
type ComponentFirmwareVersion struct {
	UUID          uuid.UUID `json:"uuid"`
	Vendor        string    `json:"vendor"`
	Model         string    `json:"model"`
	Filename      string    `json:"filename"`
	Version       string    `json:"version"`
	Component     string    `json:"component"`
	Checksum      string    `json:"checksum"`
	UpstreamURL   string    `json:"upstreamURL"`
	RepositoryURL string    `json:"repositoryURL"`
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
	f.Checksum = dbF.Checksum
	f.UpstreamURL = dbF.UpstreamURL
	f.RepositoryURL = dbF.RepositoryURL

	return nil
}

func (f *ComponentFirmwareVersion) toDBModel() (*models.ComponentFirmwareVersion, error) {
	dbF := &models.ComponentFirmwareVersion{
		Component:     f.Component,
		Vendor:        f.Vendor,
		Model:         f.Model,
		Filename:      f.Filename,
		Version:       f.Version,
		Checksum:      f.Checksum,
		UpstreamURL:   f.UpstreamURL,
		RepositoryURL: f.RepositoryURL,
	}

	if f.UUID.String() != uuid.Nil.String() {
		dbF.ID = f.UUID.String()
	}

	return dbF, nil
}
