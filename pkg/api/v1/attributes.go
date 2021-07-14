package hollow

import (
	"encoding/json"

	"go.metalkube.net/hollow/internal/db"
	"gorm.io/datatypes"
)

// Attributes provide the ability to apply namespaced settings to an entity.
// For example hardware could have attributes in the `com.equinixmetal.api` namespace
// that represents equinix metal specific attributes that are stored in the API.
// The namespace is meant to define who owns the schema and values.
type Attributes struct {
	Namespace string          `json:"namespace"`
	Values    json.RawMessage `json:"values"`
}

func (a *Attributes) fromDBModel(dbA db.Attributes) error {
	a.Namespace = dbA.Namespace
	a.Values = json.RawMessage(dbA.Values)

	return nil
}

func (a *Attributes) toDBModel() (db.Attributes, error) {
	dbA := db.Attributes{
		Namespace: a.Namespace,
		Values:    datatypes.JSON(a.Values),
	}

	return dbA, nil
}

func convertFromDBAttributes(dbAttrs []db.Attributes) ([]Attributes, error) {
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

func convertToDBAttributes(attrs []Attributes) ([]db.Attributes, error) {
	dbAttrs := []db.Attributes{}

	for _, a := range attrs {
		dbA, err := a.toDBModel()
		if err != nil {
			return nil, err
		}

		dbAttrs = append(dbAttrs, dbA)
	}

	return dbAttrs, nil
}
