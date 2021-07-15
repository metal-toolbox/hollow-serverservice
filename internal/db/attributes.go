package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Attributes provide the ability to apply namespaced settings to an entity.
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

// AttributesFilter provides the ability to filter the returned results by a
// namespaced attribute
type AttributesFilter struct {
	Namespace        string
	Keys             []string
	EqualValue       interface{}
	LessThanValue    int
	GreaterThanValue int
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
func CreateAttributes(a *Attributes) error {
	return db.Create(a).Error
}

func (f *AttributesFilter) apply(d *gorm.DB, i int) *gorm.DB {
	joinName := fmt.Sprintf("attributes_%d", i)
	column := fmt.Sprintf("%s.values", joinName)

	joinStr := fmt.Sprintf("JOIN attributes AS %s ON %s.entity_id = hardware.id AND %s.entity_type = ?", joinName, joinName, joinName)
	d = d.Joins(joinStr, "hardware")

	// filter by the namespace
	d = d.Where(fmt.Sprintf("%s.namespace = ?", joinName), f.Namespace)

	jsonKeys := jsonValueBuilder(column, f.Keys...)

	queryArgs := make([]interface{}, len(f.Keys))
	for i, v := range f.Keys {
		queryArgs[i] = v
	}

	switch {
	case f.LessThanValue != 0:
		queryArgs = append(queryArgs, f.LessThanValue)
		d = d.Where(fmt.Sprintf("(%s)::int < ?", jsonKeys), queryArgs...)
	case f.GreaterThanValue != 0:
		queryArgs = append(queryArgs, f.GreaterThanValue)
		d = d.Where(fmt.Sprintf("(%s)::int > ?", jsonKeys), queryArgs...)
	default:
		d = d.Where(datatypes.JSONQuery(column).Equals(f.EqualValue, f.Keys...))
	}

	return d
}

func jsonValueBuilder(column string, keys ...string) string {
	r := fmt.Sprintf("json_extract_path_text(%s::json,", column)

	for i := range keys {
		if i > 0 {
			r += " , "
		}

		// the actual key is represented as a "?" so that we can let GORM handle passing
		// the value in. This helps protect against SQL injection since these strings
		// could be passed in by the user.
		r += "?"
	}

	return r + ")"
}
