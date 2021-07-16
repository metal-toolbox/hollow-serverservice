package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// VersionedAttributes represents the an attribute of an entity at a specific point in time.
// You would use this over an Attribute when you want to store historical data on what the
// previous attributes were over time.
type VersionedAttributes struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	EntityID   uuid.UUID      `gorm:"<-:create;index:idx_versioned_attributes_entity;index:idx_versioned_attributes_entity_namespace;not null;"`
	EntityType string         `gorm:"<-:create;index:idx_versioned_attributes_entity;index:idx_versioned_attributes_entity_namespace;not null;"`
	Namespace  string         `gorm:"<-:create;index;index:idx_versioned_attributes_entity_namespace;not null;"`
	Values     datatypes.JSON `gorm:"<-:create;"`
	CreatedAt  time.Time
}

// BeforeSave ensures that the BIOS config passes validation checks
func (a *VersionedAttributes) BeforeSave(tx *gorm.DB) (err error) {
	if a.Namespace == "" {
		return requiredFieldMissing("VersionedAttribute", "namespace")
	}

	return nil
}

// CreateVersionedAttributes will persist VersionedAttributes into the backend datastore
func (s *Store) CreateVersionedAttributes(a *VersionedAttributes) error {
	return s.db.Create(a).Error
}

// GetVersionedAttributes will return all the BIOSConfigs for a given Hardware UUID, the list will be sorted with the newest one
// first
func (s *Store) GetVersionedAttributes(hwUUID uuid.UUID) ([]VersionedAttributes, error) {
	var al []VersionedAttributes
	if err := s.db.Where(&VersionedAttributes{EntityType: "hardware", EntityID: hwUUID}).Order("created_at desc").Find(&al).Error; err != nil {
		return nil, err
	}

	return al, nil
}
