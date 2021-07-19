package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Attributes provide the ability to store namespaced values attached to an entity.
// For example hardware could have attributes in the `com.equinixmetal.api` namespace
// that represents equinixmetal specific attributes that are stored in the API.
// The namespace is meant to define who owns the schema and values.
type Attributes struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	EntityID   uuid.UUID `gorm:"<-:create;index:idx_attributes_entity;uniqueIndex:idx_attributes_entity_namespace;not null;"`
	EntityType string    `gorm:"<-:create;index:idx_attributes_entity;uniqueIndex:idx_attributes_entity_namespace;not null;"`
	Namespace  string    `gorm:"<-:create;index;uniqueIndex:idx_attributes_entity_namespace;not null;"`
	Values     datatypes.JSON
}

// BeforeSave ensures that the attributes passes validation checks
func (a *Attributes) BeforeSave(tx *gorm.DB) (err error) {
	if a.Namespace == "" {
		return requiredFieldMissing("attributes", "namespace")
	}

	// TODO: ensure values is valid json. We can return a cleaner error than the DB does

	return nil
}

// CreateAttributes will persist attributes into the backend datastore
func (s *Store) CreateAttributes(a *Attributes) error {
	return s.db.Create(a).Error
}
