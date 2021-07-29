package hollow

import (
	"net/url"

	"go.metalkube.net/hollow/internal/db"
)

// ServerListParams allows you to filter the results
type ServerListParams struct {
	pagination
	FacilityCode                 string `form:"facility-code"`
	ComponentListParams          []ServerComponentListParams
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
	encodeServerComponentListParams(p.ComponentListParams, q)
}

func (p *ServerListParams) dbFilter() (*db.ServerFilter, error) {
	var err error

	dbF := &db.ServerFilter{
		FacilityCode: p.FacilityCode,
	}

	dbF.AttributesFilters, err = convertToDBAttributesFilter(p.AttributeListParams)
	if err != nil {
		return nil, err
	}

	dbF.VersionedAttributesFilters, err = convertToDBAttributesFilter(p.VersionedAttributeListParams)
	if err != nil {
		return nil, err
	}

	dbF.ComponentFilters, err = convertToDBComponentFilter(p.ComponentListParams)
	if err != nil {
		return nil, err
	}

	return dbF, nil
}
