package fleetdbapi

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// OperatorType is used to control what kind of search is performed for an AttributeListParams value.
type OperatorType string

// AttributeOperatorType is used to define how one or more AttributeListParam values should be SQL queried.
type AttributeOperatorType string

const (
	// AttributeLogicalOR can be passed into a AttributeListParam to SQL select the attribute an OR clause.
	AttributeLogicalOR AttributeOperatorType = "or"
	// AttributeLogicalAND is the default attribute operator, it can be passed into a AttributeListParam to SQL select the attribute a AND clause.
	AttributeLogicalAND AttributeOperatorType = "and"
)

const (
	// OperatorEqual means the value has to match the keys exactly
	OperatorEqual OperatorType = "eq"
	// OperatorLike allows you to pass in a value with % in it and match anything like it. If your string has no % in it one will be added to the end automatically
	OperatorLike = "like"
	// OperatorGreaterThan will convert the value at the given key to an int and return results that are greater than Value
	OperatorGreaterThan = "gt"
	// OperatorLessThan will convert the value at the given key to an int and return results that are less than Value
	OperatorLessThan = "lt"
)

// AttributeListParams allow you to filter the results based on attributes
type AttributeListParams struct {
	Namespace string
	Keys      []string
	Operator  OperatorType
	Value     string
	// AttributeOperatorType is used to define how this AttributeListParam value should be SQL queried
	// this value defaults to AttributeLogicalAND.
	AttributeOperator AttributeOperatorType
}

func encodeAttributesListParams(alp []AttributeListParams, key string, q url.Values) {
	for _, ap := range alp {
		value := ap.Namespace

		if len(ap.Keys) != 0 && value != "" {
			value = fmt.Sprintf("%s~%s", value, strings.Join(ap.Keys, "."))

			if ap.Operator != "" && ap.Value != "" {
				value = fmt.Sprintf("%s~%s~%s", value, ap.Operator, ap.Value)
			}
		}

		if ap.AttributeOperator != "" {
			value += "~" + string(ap.AttributeOperator)
		}

		q.Add(key, value)
	}
}

func parseQueryAttributesListParams(c *gin.Context, key string) []AttributeListParams {
	alp := []AttributeListParams{}

	attrQueryParams := c.QueryArray(key)

	for _, p := range attrQueryParams {
		// format accepted
		// "ns~keys.dot.seperated~operation~value"
		// With attr OR operator: "ns~keys.dot.seperated~operation~value~or"
		// With attr AND operator: "ns~keys.dot.seperated~operation~value~and"
		parts := strings.Split(p, "~")

		param := AttributeListParams{
			Namespace: parts[0],
		}

		if len(parts) == 1 {
			alp = append(alp, param)
			continue
		}

		param.Keys = strings.Split(parts[1], ".")

		if len(parts) == 4 || len(parts) == 5 { // nolint
			switch o := (*OperatorType)(&parts[2]); *o {
			case OperatorEqual, OperatorLike, OperatorGreaterThan, OperatorLessThan:
				param.Operator = *o
				param.Value = parts[3]
			}

			// An attribute operator is only applicable when,
			// - Theres 5 parts in the attr param string when split on `~`.
			// - Theres multiple attribute query parameters defined.
			if len(parts) == 5 && len(attrQueryParams) > 1 {
				switch o := (*AttributeOperatorType)(&parts[4]); *o {
				case AttributeLogicalAND, AttributeLogicalOR:
					param.AttributeOperator = *o
				}
			}

			// if the like search doesn't contain any % add one at the end
			if param.Operator == OperatorLike && !strings.Contains(param.Value, "%") {
				param.Value += "%"
			}
		}

		alp = append(alp, param)
	}

	return alp
}

// queryMods converts the list params into sql conditions that can be added to
// sql queries
func (p *AttributeListParams) queryMods(tblName string) qm.QueryMod {
	nsMod := qm.Where(fmt.Sprintf("%s.namespace = ?", tblName), p.Namespace)

	values := []interface{}{}
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

		values = append(values, k)
	}

	where, values := p.setJSONBWhereClause(tblName, jsonPath, values)

	// namespace AND JSONB query as a query mod
	queryMods := []qm.QueryMod{nsMod, qm.And(where, values...)}

	// OR ( namespace AND JSONB query )
	if p.AttributeOperator == AttributeLogicalOR {
		return qm.Or2(qm.Expr(queryMods...))
	}

	// AND ( namespace AND JSONB query )
	return qm.Expr(queryMods...)
}

func (p *AttributeListParams) setJSONBWhereClause(tblName, jsonPath string, values []interface{}) (string, []interface{}) {
	where := ""

	switch p.Operator {
	case OperatorLessThan:
		values = append(values, p.Value)
		where = fmt.Sprintf("json_extract_path_text(%s.data::JSONB, %s)::int < ?", tblName, jsonPath)
	case OperatorGreaterThan:
		values = append(values, p.Value)
		where = fmt.Sprintf("json_extract_path_text(%s.data::JSONB, %s)::int > ?", tblName, jsonPath)
	case OperatorLike:
		values = append(values, p.Value)
		where = fmt.Sprintf("json_extract_path_text(%s.data::JSONB, %s) LIKE ?", tblName, jsonPath)
	case OperatorEqual:
		values = append(values, p.Value)
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

	return where, values
}
