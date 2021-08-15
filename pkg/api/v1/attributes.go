package hollow

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"

	"go.metalkube.net/hollow/internal/gormdb"
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

func (a *Attributes) fromDBModel(dbA gormdb.Attributes) error {
	a.Namespace = dbA.Namespace
	a.Data = json.RawMessage(dbA.Data)
	a.CreatedAt = dbA.CreatedAt
	a.UpdatedAt = dbA.UpdatedAt

	return nil
}

func (a *Attributes) toDBModel() (gormdb.Attributes, error) {
	dbA := gormdb.Attributes{
		Namespace: a.Namespace,
		Data:      datatypes.JSON(a.Data),
	}

	return dbA, nil
}

func convertFromDBAttributes(dbAttrs []gormdb.Attributes) ([]Attributes, error) {
	attrs := []Attributes{}

	for _, dbA := range dbAttrs {
		a := Attributes{}
		if err := a.fromDBModel(dbA); err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	return attrs, nil
}

func convertToDBAttributes(attrs []Attributes) ([]gormdb.Attributes, error) {
	dbAttrs := []gormdb.Attributes{}

	for _, a := range attrs {
		dbA, err := a.toDBModel()
		if err != nil {
			return nil, err
		}

		dbAttrs = append(dbAttrs, dbA)
	}

	return dbAttrs, nil
}

func convertToDBAttributesFilter(attrs []AttributeListParams) ([]gormdb.AttributesFilter, error) {
	dbFilter := []gormdb.AttributesFilter{}

	for _, aF := range attrs {
		f := gormdb.AttributesFilter{
			Namespace:        aF.Namespace,
			Keys:             aF.Keys,
			EqualValue:       aF.EqualValue,
			LessThanValue:    aF.LessThanValue,
			GreaterThanValue: aF.GreaterThanValue,
		}
		dbFilter = append(dbFilter, f)
	}

	return dbFilter, nil
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
