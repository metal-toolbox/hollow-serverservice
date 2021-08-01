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
	ServerComponentTypeID      *uuid.UUID
	AttributesFilters          []AttributesFilter
	VersionedAttributesFilters []AttributesFilter
}

func serverComponentPreload(db *gorm.DB) *gorm.DB {
	d := db.Preload("VersionedAttributes",
		"(created_at, namespace, server_component_id) IN (?)",
		db.Table("versioned_attributes").Select("max(created_at), namespace, server_component_id").Group("namespace").Group("server_component_id"),
	)

	return d.Preload("Attributes").Preload("ServerComponentType")
}

// BeforeSave ensures that the server component type passes validation checks
func (c *ServerComponent) BeforeSave(tx *gorm.DB) (err error) {
	if c.ID.String() == uuid.Nil.String() {
		c.ID = uuid.New()
	}

	return nil
}

func (f *ServerComponentFilter) apply(d *gorm.DB) *gorm.DB {
	if f.Name != "" {
		d = d.Where("name = ?", f.Name)
	}

	if f.Vendor != "" {
		d = d.Where("vendor = ?", f.Vendor)
	}

	if f.Model != "" {
		d = d.Where("model = ?", f.Model)
	}

	if f.Serial != "" {
		d = d.Where("serial = ?", f.Serial)
	}

	if f.ServerComponentTypeID != nil {
		d = d.Where("server_component_type_id = ?", f.ServerComponentTypeID)
	}

	if f.AttributesFilters != nil {
		for i, af := range f.AttributesFilters {
			d = af.applyServerComponent(d, "server_components", i)
		}
	}

	if f.VersionedAttributesFilters != nil {
		for i, af := range f.VersionedAttributesFilters {
			d = af.applyVersionedServerComponent(d, "server_components", i)
		}
	}

	return d
}

func (f *ServerComponentFilter) nestedApply(d *gorm.DB, i int) *gorm.DB {
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

	if f.ServerComponentTypeID != nil {
		d = d.Where(joinName+".server_component_type_id = ?", f.ServerComponentTypeID)
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

// GetComponentsByServerUUID will return all the server components for a given server UUID
func (s *Store) GetComponentsByServerUUID(u uuid.UUID, filter *ServerComponentFilter, pager *Pagination) ([]ServerComponent, int64, error) {
	// if server uuid is unknown return NotFound
	if !s.ServerExists(u) {
		return nil, 0, ErrNotFound
	}

	var (
		comps []ServerComponent
		count int64
	)

	d := serverComponentPreload(s.db)

	if pager == nil {
		pager = &Pagination{}
	}

	if filter != nil {
		d = filter.apply(d)
	}

	d = d.Scopes(paginate(*pager))

	if err := d.Where("server_id = ?", u).Find(&comps).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return comps, count, nil
}
