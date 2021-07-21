package hollow

import (
	"net/url"

	"go.metalkube.net/hollow/internal/db"
)

// ServerListParams allows you to filter the results
type ServerListParams struct {
	pagination
	FacilityCode                 string `form:"facility-code"`
	AttributeListParams          []AttributeListParams
	VersionedAttributeListParams []AttributeListParams
}

func (p *ServerListParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	if p.FacilityCode != "" {
		q.Set("facility-code", p.FacilityCode)
	}

	encodeAttributesListParams(p.AttributeListParams, "attr", q)
	encodeAttributesListParams(p.VersionedAttributeListParams, "ver_attr", q)
}

func (p *ServerListParams) dbFilter() *db.ServerFilter {
	dbF := &db.ServerFilter{
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
