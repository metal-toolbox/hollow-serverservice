package db

import "gorm.io/gorm"

// HardwareFilter provides the ability to filter to hardware that is returned for
// a query
type HardwareFilter struct {
	FacilityCode               string
	AttributesFilters          []AttributesFilter
	VersionedAttributesFilters []AttributesFilter
}

func (f *HardwareFilter) apply(d *gorm.DB) *gorm.DB {
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

	return d
}
