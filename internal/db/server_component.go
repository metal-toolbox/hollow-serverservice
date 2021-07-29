package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServerComponent represents a component of a server. These can be things like
// processors, NICs, hard drives, etc.
type ServerComponent struct {
	ID                    uuid.UUID
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Name                  string
	Vendor                string
	Model                 string
	Serial                string
	ServerComponentTypeID uuid.UUID
	ServerComponentType   ServerComponentType
	ServerID              uuid.UUID
	Server                Server
	Attributes            []Attributes
	VersionedAttributes   []VersionedAttributes
}

// ServerComponentFilter provides the ability to filter the returned results by
// server components
type ServerComponentFilter struct {
	Name                       string
	Vendor                     string
	Model                      string
	Serial                     string
	ServerComponentTypeID      uuid.UUID
	AttributesFilters          []AttributesFilter
	VersionedAttributesFilters []AttributesFilter
}

// BeforeSave ensures that the server component type passes validation checks
func (c *ServerComponent) BeforeSave(tx *gorm.DB) (err error) {
	if c.ID.String() == uuid.Nil.String() {
		c.ID = uuid.New()
	}

	return nil
}

func (f *ServerComponentFilter) apply(d *gorm.DB, i int) *gorm.DB {
	joinName := fmt.Sprintf("sc_%d", i)
	joinStr := fmt.Sprintf("JOIN server_components AS %s ON %s.server_id = servers.id", joinName, joinName)

	d = d.Joins(joinStr)

	if f.Name != "" {
		d = d.Where(joinName+".name = ?", f.Name)
	}

	if f.Vendor != "" {
		d = d.Where(joinName+".vendor = ?", f.Vendor)
	}

	if f.Model != "" {
		d = d.Where(joinName+".model = ?", f.Model)
	}

	if f.Serial != "" {
		d = d.Where(joinName+".serial = ?", f.Serial)
	}

	if f.ServerComponentTypeID.String() != uuid.Nil.String() {
		d = d.Where(joinName+".server_component_type_id = ?", f.Name)
	}

	if f.AttributesFilters != nil {
		for i, af := range f.AttributesFilters {
			d = af.applyServerComponent(d, joinName, i)
		}
	}

	if f.VersionedAttributesFilters != nil {
		for i, af := range f.VersionedAttributesFilters {
			d = af.applyVersionedServerComponent(d, joinName, i)
		}
	}

	return d
}
