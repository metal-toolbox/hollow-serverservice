package db

import (
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AttributesFilter provides the ability to filter the returned results by a
// namespaced attribute
type AttributesFilter struct {
	Namespace        string
	Keys             []string
	EqualValue       interface{}
	LessThanValue    int
	GreaterThanValue int
}

func (f *AttributesFilter) applyVersioned(d *gorm.DB, i int) *gorm.DB {
	joinName := fmt.Sprintf("ver_attr_%d", i)
	joinStr := fmt.Sprintf("JOIN versioned_attributes AS %s ON %s.entity_id = servers.id AND %s.entity_type = ? AND %s.created_at=(select max(created_at) from versioned_attributes where entity_id = servers.id AND entity_type = ? AND namespace = ?)", joinName, joinName, joinName, joinName)

	d = d.Joins(joinStr, "servers", "servers", f.Namespace)

	return f.addFilter(d, joinName)
}

func (f *AttributesFilter) apply(d *gorm.DB, i int) *gorm.DB {
	joinName := fmt.Sprintf("attr_%d", i)
	joinStr := fmt.Sprintf("JOIN attributes AS %s ON %s.entity_id = servers.id AND %s.entity_type = ?", joinName, joinName, joinName)

	d = d.Joins(joinStr, "servers", "servers")

	return f.addFilter(d, joinName)
}

func (f *AttributesFilter) addFilter(d *gorm.DB, joinName string) *gorm.DB {
	column := fmt.Sprintf("%s.values", joinName)

	// filter by the namespace
	d = d.Where(fmt.Sprintf("%s.namespace = ?", joinName), f.Namespace)

	jsonKeys := jsonValueBuilder(column, f.Keys...)

	queryArgs := make([]interface{}, len(f.Keys))
	for i, v := range f.Keys {
		queryArgs[i] = v
	}

	switch {
	case f.LessThanValue != 0:
		queryArgs = append(queryArgs, f.LessThanValue)
		d = d.Where(fmt.Sprintf("(%s)::int < ?", jsonKeys), queryArgs...)
	case f.GreaterThanValue != 0:
		queryArgs = append(queryArgs, f.GreaterThanValue)
		d = d.Where(fmt.Sprintf("(%s)::int > ?", jsonKeys), queryArgs...)
	default:
		d = d.Where(datatypes.JSONQuery(column).Equals(f.EqualValue, f.Keys...))
	}

	return d
}

func jsonValueBuilder(column string, keys ...string) string {
	r := fmt.Sprintf("json_extract_path_text(%s::json,", column)

	for i := range keys {
		if i > 0 {
			r += " , "
		}

		// the actual key is represented as a "?" so that we can let GORM handle passing
		// the value in. This helps protect against SQL injection since these strings
		// could be passed in by the user.
		r += "?"
	}

	return r + ")"
}
