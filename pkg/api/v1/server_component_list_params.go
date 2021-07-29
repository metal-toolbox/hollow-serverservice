package hollow

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// ServerComponentListParams allows you to filter the results by server components
type ServerComponentListParams struct {
	Name                         string
	Vendor                       string
	Model                        string
	Serial                       string
	ServerComponentTypeUUID      uuid.UUID
	AttributeListParams          []AttributeListParams
	VersionedAttributeListParams []AttributeListParams
}

func (p *ServerComponentListParams) empty() bool {
	switch {
	case p.Name != "",
		p.Vendor != "",
		p.Model != "",
		p.Serial != "",
		p.ServerComponentTypeUUID.String() != uuid.Nil.String(),
		len(p.AttributeListParams) != 0,
		len(p.VersionedAttributeListParams) != 0:
		return false
	default:
		return true
	}
}

func convertToDBComponentFilter(sclp []ServerComponentListParams) ([]db.ServerComponentFilter, error) {
	var err error

	dbFilters := []db.ServerComponentFilter{}

	for _, p := range sclp {
		dbF := db.ServerComponentFilter{
			Name:                  p.Name,
			Vendor:                p.Vendor,
			Model:                 p.Model,
			Serial:                p.Serial,
			ServerComponentTypeID: p.ServerComponentTypeUUID,
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
		keyPrefix := fmt.Sprintf("sc_%d_", i)

		if sp.Name != "" {
			q.Set(keyPrefix+"name", sp.Name)
		}

		if sp.Vendor != "" {
			q.Set(keyPrefix+"vendor", sp.Vendor)
		}

		if sp.Model != "" {
			q.Set(keyPrefix+"model", sp.Model)
		}

		if sp.Serial != "" {
			q.Set(keyPrefix+"serial", sp.Serial)
		}

		if sp.ServerComponentTypeUUID.String() != uuid.Nil.String() {
			q.Set(keyPrefix+"server_component_type_uuid", sp.Name)
		}

		encodeAttributesListParams(sp.AttributeListParams, keyPrefix+"attr", q)
		encodeAttributesListParams(sp.VersionedAttributeListParams, keyPrefix+"ver_attr", q)
	}
}

func parseQueryServerComponentsListParams(c *gin.Context) ([]ServerComponentListParams, error) {
	var err error

	sclp := []ServerComponentListParams{}
	i := 0

	for {
		keyPrefix := fmt.Sprintf("sc_%d_", i)

		var u uuid.UUID

		if c.Query(keyPrefix+"server_component_type_uuid") != "" {
			u, err = uuid.Parse(c.Query(keyPrefix + "server_component_type_uuid"))
			if err != nil {
				return nil, err
			}
		}

		p := ServerComponentListParams{
			Name:                    c.Query(keyPrefix + "name"),
			Vendor:                  c.Query(keyPrefix + "vendor"),
			Model:                   c.Query(keyPrefix + "model"),
			Serial:                  c.Query(keyPrefix + "serial"),
			ServerComponentTypeUUID: u,
		}

		alp, err := parseQueryAttributesListParams(c, keyPrefix+"attr")
		if err != nil {
			return nil, err
		}

		p.AttributeListParams = alp

		valp, err := parseQueryAttributesListParams(c, keyPrefix+"ver_attr")
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
