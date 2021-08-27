package dcim

import (
	"encoding/json"
	"time"

	"github.com/volatiletech/sqlboiler/v4/types"

	"go.hollow.sh/serverservice/internal/models"
)

// VersionedAttributes represents a set of attributes of an entity at a given time
type VersionedAttributes struct {
	Namespace      string          `json:"namespace" binding:"required"`
	Data           json.RawMessage `json:"data" binding:"required"`
	Tally          int             `json:"tally"`
	LastReportedAt time.Time       `json:"last_reported_at"`
	CreatedAt      time.Time       `json:"created_at"`
}

func (a *VersionedAttributes) toDBModel() *models.VersionedAttribute {
	return &models.VersionedAttribute{
		Namespace: a.Namespace,
		Data:      types.JSON(a.Data),
	}
}

func (a *VersionedAttributes) fromDBModel(dba *models.VersionedAttribute) error {
	a.CreatedAt = dba.CreatedAt.Time
	a.LastReportedAt = dba.UpdatedAt.Time
	a.Tally = int(dba.Tally)
	a.Namespace = dba.Namespace
	a.Data = json.RawMessage(dba.Data)

	return nil
}

func convertFromDBVersionedAttributes(dbAttrs models.VersionedAttributeSlice) ([]VersionedAttributes, error) {
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
