package serverservice

import (
	"time"

	"github.com/google/uuid"

	"go.hollow.sh/serverservice/internal/models"
)

// ComponentFirmwareVersion represents a firmware file
type ComponentFirmwareVersion struct {
	UUID          uuid.UUID `json:"uuid"`
	Vendor        string    `json:"vendor" binding:"required,lowercase"`
	Model         string    `json:"model" binding:"required,lowercase"`
	Filename      string    `json:"filename" binding:"required"`
	Version       string    `json:"version" binding:"required"`
	Component     string    `json:"component" binding:"required,lowercase"`
	Checksum      string    `json:"checksum" binding:"required,lowercase"`
	UpstreamURL   string    `json:"upstreamURL" binding:"required"`
	RepositoryURL string    `json:"repositoryURL" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
	f.CreatedAt = dbF.CreatedAt.Time
	f.UpdatedAt = dbF.UpdatedAt.Time

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
