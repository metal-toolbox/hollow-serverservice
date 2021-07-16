package hollow

import (
	"net/url"

	"go.metalkube.net/hollow/internal/db"
)

// HardwareListParams allows you to filter the results
type HardwareListParams struct {
	FacilityCode                 string `form:"facility-code"`
	AttributeListParams          []AttributeListParams
	VersionedAttributeListParams []AttributeListParams
}

func (p *HardwareListParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	if p.FacilityCode != "" {
		q.Set("facility-code", p.FacilityCode)
	}

	encodeAttributesListParams(p.AttributeListParams, "attr", q)
	encodeAttributesListParams(p.VersionedAttributeListParams, "ver_attr", q)
}

func (p *HardwareListParams) dbFilter() *db.HardwareFilter {
	dbF := &db.HardwareFilter{
		FacilityCode: p.FacilityCode,
	}

	for _, aF := range p.AttributeListParams {
		a := db.AttributesFilter{
			Namespace:        aF.Namespace,
			Keys:             aF.Keys,
			EqualValue:       aF.EqualValue,
			LessThanValue:    aF.LessThanValue,
			GreaterThanValue: aF.GreaterThanValue,
		}
		dbF.AttributesFilters = append(dbF.AttributesFilters, a)
	}

	for _, aF := range p.VersionedAttributeListParams {
		a := db.AttributesFilter{
			Namespace:        aF.Namespace,
			Keys:             aF.Keys,
			EqualValue:       aF.EqualValue,
			LessThanValue:    aF.LessThanValue,
			GreaterThanValue: aF.GreaterThanValue,
		}
		dbF.VersionedAttributesFilters = append(dbF.VersionedAttributesFilters, a)
	}

	return dbF
}
