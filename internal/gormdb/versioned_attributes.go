package gormdb

import (
	"encoding/json"
	"reflect"
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
	ID                uuid.UUID
	ServerID          *uuid.UUID `gorm:"<-:create;"`
	Server            *Server
	ServerComponentID *uuid.UUID `gorm:"<-:create;"`
	ServerComponent   *ServerComponent
	Namespace         string         `gorm:"<-:create;"`
	Data              datatypes.JSON `gorm:"<-:create;"`
	Tally             int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// BeforeSave ensures that the VersionedAttributes passes validation checks
func (a *VersionedAttributes) BeforeSave(tx *gorm.DB) (err error) {
	if a.ID.String() == uuid.Nil.String() {
		a.ID = uuid.New()
	}

	if a.Namespace == "" {
		return requiredFieldMissing("VersionedAttributes", "namespace")
	}

	return nil
}

// CreateVersionedAttributes will persist VersionedAttributes into the backend datastore
func (s *Store) CreateVersionedAttributes(entity interface{}, a *VersionedAttributes) error {
	var existing VersionedAttributes

	if err := a.BeforeSave(s.db); err != nil {
		return err
	}

	d := s.db.Model(entity).Where("namespace = ?", a.Namespace).Order("created_at desc").Limit(1).Association("VersionedAttributes")

	err := d.Find(&existing)
	if err != nil {
		return err
	}

	if existing.ID.String() != uuid.Nil.String() && areEqualJSON(json.RawMessage(existing.Data), json.RawMessage(a.Data)) {
		existing.Tally++
		return s.db.Updates(&existing).Error
	}

	return s.db.Model(entity).Association("VersionedAttributes").Append(a)
}

// ListVersionedAttributes will return all the VersionedAttributes for a given server UUID, the list will be sorted with the newest one
// first
func (s *Store) ListVersionedAttributes(srvUUID uuid.UUID, pager *Pagination) ([]VersionedAttributes, int64, error) {
	var (
		al    []VersionedAttributes
		count int64
	)

	if pager == nil {
		pager = &Pagination{}
	}

	d := s.db.Where(&VersionedAttributes{ServerID: &srvUUID})

	if err := d.Scopes(paginate(*pager)).Order("created_at desc").Find(&al).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return al, count, nil
}

// GetVersionedAttributes will return all the VersionedAttributes for a given server UUID and namespace, the list will be sorted with the newest one
// first
func (s *Store) GetVersionedAttributes(srvUUID uuid.UUID, ns string, pager *Pagination) ([]VersionedAttributes, int64, error) {
	var (
		al    []VersionedAttributes
		count int64
	)

	if pager == nil {
		pager = &Pagination{}
	}

	d := s.db.Where(&VersionedAttributes{ServerID: &srvUUID, Namespace: ns})

	if err := d.Scopes(paginate(*pager)).Order("created_at desc").Find(&al).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return al, count, nil
}

func areEqualJSON(s1, s2 json.RawMessage) bool {
	var (
		o1 interface{}
		o2 interface{}
	)

	if err := json.Unmarshal([]byte(s1), &o1); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(s2), &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
