package hollow

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"

	"go.metalkube.net/hollow/internal/db"
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

func (p *ServerComponentListParams) empty() bool {
	switch {
	case p.Name != "",
		p.Vendor != "",
		p.Model != "",
		p.Serial != "",
		p.ServerComponentType != "",
		len(p.AttributeListParams) != 0,
		len(p.VersionedAttributeListParams) != 0:
		return false
	default:
		return true
	}
}

func convertToDBComponentFilter(r *Router, sclp []ServerComponentListParams) ([]db.ServerComponentFilter, error) {
	var err error

	dbFilters := []db.ServerComponentFilter{}

	for _, p := range sclp {
		dbF := db.ServerComponentFilter{
			Name:   p.Name,
			Vendor: p.Vendor,
			Model:  p.Model,
			Serial: p.Serial,
		}

		if p.ServerComponentType != "" {
			fmt.Printf("\n\n\nLooking Up Slug: %s\n\n", p.ServerComponentType)

			sct, err := r.Store.FindServerComponentTypeBySlug(p.ServerComponentType)
			if err != nil {
				return nil, err
			}

			fmt.Printf("Found Type By Slug, Setting ID to: %s\n\n", sct.ID)

			dbF.ServerComponentTypeID = &sct.ID
		}

		dbF.AttributesFilters, err = convertToDBAttributesFilter(p.AttributeListParams)
		if err != nil {
			return nil, err
		}

		dbF.VersionedAttributesFilters, err = convertToDBAttributesFilter(p.VersionedAttributeListParams)
		if err != nil {
			return nil, err
		}

		dbFilters = append(dbFilters, dbF)
	}

	return dbFilters, nil
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

		p := ServerComponentListParams{
			Name:                queryMap["name"],
			Vendor:              queryMap["vendor"],
			Model:               queryMap["model"],
			Serial:              queryMap["serial"],
			ServerComponentType: queryMap["type"],
		}

		alp, err := parseQueryAttributesListParams(c, keyPrefix+"_attr")
		if err != nil {
			return nil, err
		}

		p.AttributeListParams = alp

		valp, err := parseQueryAttributesListParams(c, keyPrefix+"_ver_attr")
		if err != nil {
			return nil, err
		}

		p.VersionedAttributeListParams = valp

		if p.empty() {
			// if no attributes are set then one wasn't passed in. Break out of the loop
			break
		}

		sclp = append(sclp, p)
		i++
	}

	return sclp, nil
}
