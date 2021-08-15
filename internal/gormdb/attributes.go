package gormdb

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Attributes provide the ability to store namespaced values attached to an entity.
// For example servers could have attributes in the `com.equinixmetal.api` namespace
// that represents equinixmetal specific attributes that are stored in the API.
// The namespace is meant to define who owns the schema and values.
type Attributes struct {
	ID                uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ServerID          *uuid.UUID `gorm:"<-:create;"`
	Server            *Server
	ServerComponentID *uuid.UUID `gorm:"<-:create;"`
	ServerComponent   *ServerComponent
	Namespace         string `gorm:"<-:create;"`
	Data              datatypes.JSON
}

// BeforeSave ensures that the attributes passes validation checks
func (a *Attributes) BeforeSave(tx *gorm.DB) (err error) {
	if a.ID.String() == uuid.Nil.String() {
		a.ID = uuid.New()
	}

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

// DeleteAttributes will persist attributes into the backend datastore
func (s *Store) DeleteAttributes(a *Attributes) error {
	return s.db.Delete(a).Error
}

// GetAttributesByServerUUID will return all the attributes for a given server UUID
func (s *Store) GetAttributesByServerUUID(u uuid.UUID, pager *Pagination) ([]Attributes, int64, error) {
	// if server uuid is unknown return NotFound
	if !s.ServerExists(u) {
		return nil, 0, ErrNotFound
	}

	var (
		attrs []Attributes
		count int64
	)

	if pager == nil {
		pager = &Pagination{}
	}

	d := s.db.Preload(clause.Associations).Scopes(paginate(*pager))

	if err := d.Where("server_id = ?", u).Find(&attrs).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return attrs, count, nil
}

// GetAttributesByServerUUIDAndNamespace will return attributes for a given server UUID and namespace
func (s *Store) GetAttributesByServerUUIDAndNamespace(u uuid.UUID, ns string) (*Attributes, error) {
	var attr Attributes

	d := s.db.Preload(clause.Associations)

	if err := d.Where("server_id = ?", u).Where("namespace = ?", ns).First(&attr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &attr, nil
}

// FindOrInitAttributesByServerUUIDAndNamespace will return attributes for a given server UUID and namespace. If one
// doesn't exist an empty one will be returned which can be saved
func (s *Store) FindOrInitAttributesByServerUUIDAndNamespace(u uuid.UUID, ns string) (*Attributes, error) {
	if !s.ServerExists(u) || ns == "" {
		return nil, ErrNotFound
	}

	var attr Attributes

	d := s.db.Preload(clause.Associations)

	if err := d.FirstOrInit(&attr, Attributes{ServerID: &u, Namespace: ns}).Error; err != nil {
		return nil, err
	}

	return &attr, nil
}

// UpdateAttributesByServerUUIDAndNamespace allow you to update the data stored in a given namespace for a server
func (s *Store) UpdateAttributesByServerUUIDAndNamespace(u uuid.UUID, ns string, data json.RawMessage) error {
	attr, err := s.FindOrInitAttributesByServerUUIDAndNamespace(u, ns)
	if err != nil {
		return err
	}

	attr.Data = datatypes.JSON(data)

	return s.db.Save(&attr).Error
}
