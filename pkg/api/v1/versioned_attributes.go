package hollow

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"

	"go.metalkube.net/hollow/internal/db"
)

// VersionedAttributes represents a set of attributes of an entity at a given time
type VersionedAttributes struct {
	Namespace string          `json:"namespace" binding:"required"`
	Values    json.RawMessage `json:"values" binding:"required"`
	CreatedAt time.Time       `json:"created_at"`
}

func (a *VersionedAttributes) toDBModel() (*db.VersionedAttributes, error) {
	dbc := &db.VersionedAttributes{
		Namespace: a.Namespace,
		Values:    datatypes.JSON(a.Values),
	}

	return dbc, nil
}

func (a *VersionedAttributes) fromDBModel(dba db.VersionedAttributes) error {
	a.CreatedAt = dba.CreatedAt
	a.Namespace = dba.Namespace
	a.Values = json.RawMessage(dba.Values)

	return nil
}

func convertToDBVersionedAttributes(attrs []VersionedAttributes) ([]db.VersionedAttributes, error) {
	dbVerAttrs := []db.VersionedAttributes{}

	for _, a := range attrs {
		dbVA, err := a.toDBModel()
		if err != nil {
			return nil, err
		}

		dbVerAttrs = append(dbVerAttrs, *dbVA)
	}

	return dbVerAttrs, nil
}

func convertFromDBVersionedAttributes(dbAttrs []db.VersionedAttributes) ([]VersionedAttributes, error) {
	attrs := []VersionedAttributes{}

	for _, dbA := range dbAttrs {
		a := VersionedAttributes{}
		if err := a.fromDBModel(dbA); err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	return attrs, nil
}
