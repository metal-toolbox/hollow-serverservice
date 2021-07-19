package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// VersionedAttributes provide the ability to store namespaced values attached
// to an entity, tied to a specific timestamp. You would use VersionedAttributes
// over Attributes when you want to store historical data on what the previous
// values were.
type VersionedAttributes struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	EntityID   uuid.UUID      `gorm:"<-:create;index:idx_versioned_attributes_entity;index:idx_versioned_attributes_entity_namespace;not null;"`
	EntityType string         `gorm:"<-:create;index:idx_versioned_attributes_entity;index:idx_versioned_attributes_entity_namespace;not null;"`
	Namespace  string         `gorm:"<-:create;index;index:idx_versioned_attributes_entity_namespace;not null;"`
	Values     datatypes.JSON `gorm:"<-:create;"`
	CreatedAt  time.Time
}

// BeforeSave ensures that the VersionedAttributes passes validation checks
func (a *VersionedAttributes) BeforeSave(tx *gorm.DB) (err error) {
	if a.Namespace == "" {
		return requiredFieldMissing("VersionedAttributes", "namespace")
	}

	return nil
}

// CreateVersionedAttributes will persist VersionedAttributes into the backend datastore
func (s *Store) CreateVersionedAttributes(entity interface{}, a *VersionedAttributes) error {
	// return s.db.Create(a).Error
	return s.db.Model(entity).Association("VersionedAttributes").Append(a)
}

// GetVersionedAttributes will return all the VersionedAttributes for a given Hardware UUID, the list will be sorted with the newest one
// first
func (s *Store) GetVersionedAttributes(hwUUID uuid.UUID) ([]VersionedAttributes, error) {
	var al []VersionedAttributes
	if err := s.db.Where(&VersionedAttributes{EntityType: "hardware", EntityID: hwUUID}).Order("created_at desc").Find(&al).Error; err != nil {
		return nil, err
	}

	return al, nil
}
