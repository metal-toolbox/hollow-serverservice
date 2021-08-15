package hollow

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"

	"go.metalkube.net/hollow/internal/gormdb"
)

// VersionedAttributes represents a set of attributes of an entity at a given time
type VersionedAttributes struct {
	Namespace      string          `json:"namespace" binding:"required"`
	Data           json.RawMessage `json:"data" binding:"required"`
	Tally          int             `json:"tally"`
	LastReportedAt time.Time       `json:"last_reported_at"`
	CreatedAt      time.Time       `json:"created_at"`
}

func (a *VersionedAttributes) toDBModel() (*gormdb.VersionedAttributes, error) {
	dbc := &gormdb.VersionedAttributes{
		Namespace: a.Namespace,
		Data:      datatypes.JSON(a.Data),
	}

	return dbc, nil
}

func (a *VersionedAttributes) fromDBModel(dba gormdb.VersionedAttributes) error {
	a.CreatedAt = dba.CreatedAt
	a.LastReportedAt = dba.UpdatedAt
	a.Tally = dba.Tally
	a.Namespace = dba.Namespace
	a.Data = json.RawMessage(dba.Data)

	return nil
}

func convertToDBVersionedAttributes(attrs []VersionedAttributes) ([]gormdb.VersionedAttributes, error) {
	dbVerAttrs := []gormdb.VersionedAttributes{}

	for _, a := range attrs {
		dbVA, err := a.toDBModel()
		if err != nil {
			return nil, err
		}

		dbVerAttrs = append(dbVerAttrs, *dbVA)
	}

	return dbVerAttrs, nil
}

func convertFromDBVersionedAttributes(dbAttrs []gormdb.VersionedAttributes) ([]VersionedAttributes, error) {
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
