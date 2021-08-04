package db

import "gorm.io/gorm"

// ServerFilter provides the ability to filter the list of servers that are
// returned for a query
type ServerFilter struct {
	FacilityCode               string
	ComponentFilters           []ServerComponentFilter
	AttributesFilters          []AttributesFilter
	VersionedAttributesFilters []AttributesFilter
}

func (f *ServerFilter) apply(d *gorm.DB) *gorm.DB {
	if f.FacilityCode != "" {
		d = d.Where("facility_code = ?", f.FacilityCode)
	}

	if f.AttributesFilters != nil {
		for i, af := range f.AttributesFilters {
			d = af.apply(d, i)
		}
	}

	if f.VersionedAttributesFilters != nil {
		for i, af := range f.VersionedAttributesFilters {
			d = af.applyVersioned(d, i)
		}
	}

	if f.ComponentFilters != nil {
		for i, cf := range f.ComponentFilters {
			d = cf.nestedApply(d, i)
		}
	}

	return d
}
