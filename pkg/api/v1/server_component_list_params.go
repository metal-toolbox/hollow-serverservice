package fleetdbapi

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// ServerComponentListParams allows you to filter the results by server components
type ServerComponentListParams struct {
	Name                         string
	Vendor                       string
	Model                        string
	Serial                       string
	ServerComponentType          string
	AttributeListParams          []AttributeListParams
	VersionedAttributeListParams []AttributeListParams
	Pagination                   *PaginationParams
}

// setQuery implements the queryParams interface
func (p *ServerComponentListParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	encodeAttributesListParams(p.AttributeListParams, "attr", q)
	encodeAttributesListParams(p.VersionedAttributeListParams, "ver_attr", q)
	encodeServerComponentListParams([]ServerComponentListParams{*p}, q)
	p.Pagination.setQuery(q)
}

// queryMods converts the list params into sql conditions that can be added to
// sql queries
func (p *ServerComponentListParams) queryMods(tblName string) qm.QueryMod {
	mods := []qm.QueryMod{}

	if p.Name != "" {
		mods = append(mods, qm.Where(fmt.Sprintf("%s.name = ?", tblName), p.Name))
	}

	if p.Vendor != "" {
		mods = append(mods, qm.Where(fmt.Sprintf("%s.vendor = ?", tblName), p.Vendor))
	}

	if p.Model != "" {
		mods = append(mods, qm.Where(fmt.Sprintf("%s.model = ?", tblName), p.Model))
	}

	if p.Serial != "" {
		mods = append(mods, qm.Where(fmt.Sprintf("%s.serial = ?", tblName), p.Serial))
	}

	if p.ServerComponentType != "" {
		joinTblName := fmt.Sprintf("%s_sct", tblName)
		whereStmt := fmt.Sprintf("server_component_types as %s on %s.server_component_type_id = %s.id", joinTblName, tblName, joinTblName)
		mods = append(mods, qm.LeftOuterJoin(whereStmt))
		mods = append(mods, qm.Where(fmt.Sprintf("%s.slug = ?", joinTblName), p.ServerComponentType))
	}

	for i, lp := range p.AttributeListParams {
		tableName := fmt.Sprintf("%s_attr_%d", tblName, i)
		whereStmt := fmt.Sprintf("attributes as %s on %s.server_component_id = %s.id", tableName, tableName, tblName)
		mods = append(mods, qm.LeftOuterJoin(whereStmt))

		mods = append(mods, lp.queryMods(tableName))
	}

	for i, lp := range p.VersionedAttributeListParams {
		tableName := fmt.Sprintf("%s_ver_attr_%d", tblName, i)
		whereStmt := fmt.Sprintf("versioned_attributes as %s on %s.server_component_id = %s.id AND %s.created_at=(select max(created_at) from versioned_attributes where server_component_id = %s.id AND namespace = ?)", tableName, tableName, tblName, tableName, tblName)
		mods = append(mods, qm.LeftOuterJoin(whereStmt, lp.Namespace))
		mods = append(mods, lp.queryMods(tableName))
	}

	return qm.Expr(mods...)
}

func encodeServerComponentListParams(sclp []ServerComponentListParams, q url.Values) {
	for i, sp := range sclp {
		keyPrefix := fmt.Sprintf("sc_%d", i)

		if sp.Name != "" {
			q.Set(keyPrefix+"[name]", sp.Name)
		}

		if sp.Vendor != "" {
			q.Set(keyPrefix+"[vendor]", sp.Vendor)
		}

		if sp.Model != "" {
			q.Set(keyPrefix+"[model]", sp.Model)
		}

		if sp.Serial != "" {
			q.Set(keyPrefix+"[serial]", sp.Serial)
		}

		if sp.ServerComponentType != "" {
			q.Set(keyPrefix+"[type]", sp.ServerComponentType)
		}

		encodeAttributesListParams(sp.AttributeListParams, keyPrefix+"_attr", q)
		encodeAttributesListParams(sp.VersionedAttributeListParams, keyPrefix+"_ver_attr", q)
	}
}

func parseQueryServerComponentsListParams(c *gin.Context) ([]ServerComponentListParams, error) {
	sclp := []ServerComponentListParams{}
	i := 0

	for {
		keyPrefix := fmt.Sprintf("sc_%d", i)

		queryMap := c.QueryMap(keyPrefix)

		aListParams := parseQueryAttributesListParams(c, keyPrefix+"_attr")
		vaListParams := parseQueryAttributesListParams(c, keyPrefix+"_ver_attr")

		// no parameters were passed in, break out of loop
		if len(queryMap) == 0 && len(aListParams) == 0 && len(vaListParams) == 0 {
			break
		}

		p := ServerComponentListParams{
			Name:                queryMap["name"],
			Vendor:              queryMap["vendor"],
			Model:               queryMap["model"],
			Serial:              queryMap["serial"],
			ServerComponentType: queryMap["type"],
		}

		if len(aListParams) > 0 {
			p.AttributeListParams = aListParams
		}

		if len(vaListParams) > 0 {
			p.VersionedAttributeListParams = vaListParams
		}

		sclp = append(sclp, p)
		i++
	}

	return sclp, nil
}
