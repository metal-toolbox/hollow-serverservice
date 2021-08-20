package dcim

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"

	"go.hollow.sh/dcim/internal/models"
)

// Attributes provide the ability to apply namespaced settings to an entity.
// For example servers could have attributes in the `com.equinixmetal.api` namespace
// that represents equinix metal specific attributes that are stored in the API.
// The namespace is meant to define who owns the schema and values.
type Attributes struct {
	Namespace string          `json:"namespace"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// AttributeListParams allow you to filter the results based on attributes
type AttributeListParams struct {
	Namespace        string   `form:"namespace" query:"namespace"`
	Keys             []string `form:"keys" query:"keys"`
	EqualValue       string   `form:"equals" query:"equals"`
	LessThanValue    int      `form:"less-than" query:"less-than"`
	GreaterThanValue int      `form:"greater-than" query:"greater-than"`
}

func (a *Attributes) fromDBModel(dbA *models.Attribute) error {
	a.Namespace = dbA.Namespace
	a.Data = json.RawMessage(dbA.Data)
	a.CreatedAt = dbA.CreatedAt.Time
	a.UpdatedAt = dbA.UpdatedAt.Time

	return nil
}

func (a *Attributes) toDBModel() (*models.Attribute, error) {
	dbA := &models.Attribute{
		Namespace: a.Namespace,
		Data:      types.JSON(a.Data),
	}

	return dbA, nil
}

func convertFromDBAttributes(dbAttrs models.AttributeSlice) ([]Attributes, error) {
	attrs := []Attributes{}
	if dbAttrs == nil {
		return attrs, nil
	}

	for _, dbA := range dbAttrs {
		a := Attributes{}
		if err := a.fromDBModel(dbA); err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	return attrs, nil
}

func encodeAttributesListParams(alp []AttributeListParams, key string, q url.Values) {
	for _, ap := range alp {
		value := ap.Namespace

		if len(ap.Keys) != 0 {
			value = fmt.Sprintf("%s~%s", value, strings.Join(ap.Keys, "."))

			switch {
			case ap.LessThanValue != 0:
				value = fmt.Sprintf("%s~lt~%d", value, ap.LessThanValue)
			case ap.GreaterThanValue != 0:
				value = fmt.Sprintf("%s~gt~%d", value, ap.GreaterThanValue)
			case ap.EqualValue != "":
				value = fmt.Sprintf("%s~eq~%s", value, ap.EqualValue)
			}
		}

		q.Add(key, value)
	}
}

func parseQueryAttributesListParams(c *gin.Context, key string) ([]AttributeListParams, error) {
	var err error

	alp := []AttributeListParams{}

	for _, p := range c.QueryArray(key) {
		// format is "ns~keys.dot.seperated~operation~value"
		parts := strings.Split(p, "~")

		param := AttributeListParams{
			Namespace: parts[0],
		}

		if len(parts) == 1 {
			alp = append(alp, param)
			continue
		}

		param.Keys = strings.Split(parts[1], ".")

		if len(parts) == 4 { //nolint
			switch op := parts[2]; op {
			case "lt":
				param.LessThanValue, err = strconv.Atoi(parts[3])
				if err != nil {
					return nil, err
				}
			case "gt":
				param.GreaterThanValue, err = strconv.Atoi(parts[3])
				if err != nil {
					return nil, err
				}
			case "eq":
				param.EqualValue = parts[3]
			}
		}

		alp = append(alp, param)
	}

	return alp, nil
}

// queryMods converts the list params into sql conditions that can be added to
// sql queries
func (p *AttributeListParams) queryMods(tblName string) qm.QueryMod {
	nsMod := qm.Where(fmt.Sprintf("%s.namespace = ?", tblName), p.Namespace)

	sqlValues := []interface{}{}
	jsonPath := ""

	// If we only have a namespace and no keys we are limiting by namespace only
	if len(p.Keys) == 0 {
		return nsMod
	}

	for i, k := range p.Keys {
		if i > 0 {
			jsonPath += " , "
		}
		// the actual key is represented as a "?" this helps protect against SQL
		// injection since these strings are passed in by the user.
		jsonPath += "?"

		sqlValues = append(sqlValues, k)
	}

	where := ""

	switch {
	case p.LessThanValue != 0:
		sqlValues = append(sqlValues, p.LessThanValue)
		where = fmt.Sprintf("json_extract_path_text(%s.data::JSONB, %s)::int < ?", tblName, jsonPath)
	case p.GreaterThanValue != 0:
		sqlValues = append(sqlValues, p.GreaterThanValue)
		where = fmt.Sprintf("json_extract_path_text(%s.data::JSONB, %s)::int > ?", tblName, jsonPath)
	case p.EqualValue != "":
		sqlValues = append(sqlValues, p.EqualValue)
		where = fmt.Sprintf("json_extract_path_text(%s.data::JSONB, %s) = ?", tblName, jsonPath)
	default:
		// we only have keys so we just want to ensure the key is there
		where = fmt.Sprintf("%s.data::JSONB", tblName)

		if len(p.Keys) != 0 {
			for range p.Keys[0 : len(p.Keys)-1] {
				where += " -> ?"
			}

			// query is existing_where ? key
			where += " \\? ?"
		}
	}

	return qm.Expr(
		nsMod,
		qm.And(where, sqlValues...),
	)
}
