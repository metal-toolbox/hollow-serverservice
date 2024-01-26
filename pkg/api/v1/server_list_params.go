package serverservice

import (
	"fmt"
	"net/url"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerListParams allows you to filter the results
type ServerListParams struct {
	FacilityCode                 string `form:"facility-code"`
	ComponentListParams          []ServerComponentListParams
	AttributeListParams          []AttributeListParams
	IncludeDeleted               bool `form:"include-deleted"`
	VersionedAttributeListParams []AttributeListParams
	PaginationParams             *PaginationParams
}

func (p *ServerListParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	if p.FacilityCode != "" {
		q.Set("facility-code", p.FacilityCode)
	}

	if p.IncludeDeleted {
		q.Set("include-deleted", "true")
	}

	encodeAttributesListParams(p.AttributeListParams, "attr", q)
	encodeAttributesListParams(p.VersionedAttributeListParams, "ver_attr", q)
	encodeServerComponentListParams(p.ComponentListParams, q)
	p.PaginationParams.setQuery(q)
}

// queryMods converts the list params into sql conditions that can be added to
// sql queries
func (p *ServerListParams) queryMods() []qm.QueryMod {
	mods := []qm.QueryMod{}

	if p.FacilityCode != "" {
		m := models.ServerWhere.FacilityCode.EQ(null.StringFrom(p.FacilityCode))
		mods = append(mods, m)
	}

	mods = append(mods, qm.Distinct("servers.*"))

	for i, lp := range p.AttributeListParams {
		tableName := fmt.Sprintf("attr_%d", i)
		whereStmt := fmt.Sprintf("attributes as %s on %s.server_id = servers.id", tableName, tableName)
		mods = append(mods, qm.LeftOuterJoin(whereStmt))

		mods = append(mods, lp.queryMods(tableName))
	}

	for i, lp := range p.VersionedAttributeListParams {
		tableName := fmt.Sprintf("ver_attr_%d", i)
		whereStmt := fmt.Sprintf("versioned_attributes as %s on %s.server_id = servers.id AND %s.created_at=(select max(created_at) from versioned_attributes where server_id = servers.id AND namespace = ?)", tableName, tableName, tableName)
		mods = append(mods, qm.LeftOuterJoin(whereStmt, lp.Namespace))
		mods = append(mods, lp.queryMods(tableName))
	}

	for i, lp := range p.ComponentListParams {
		tableName := fmt.Sprintf("sc_%d", i)
		whereStmt := fmt.Sprintf("server_components as %s on %s.server_id = servers.id", tableName, tableName)
		mods = append(mods, qm.LeftOuterJoin(whereStmt))
		mods = append(mods, lp.queryMods(tableName))
	}

	if p.IncludeDeleted {
		mods = append(mods, qm.WithDeleted())
	}

	return mods
}
