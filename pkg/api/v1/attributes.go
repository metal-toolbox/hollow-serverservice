package serverservice

import (
	"encoding/json"
	"time"

	"github.com/volatiletech/sqlboiler/v4/types"

	"go.hollow.sh/serverservice/internal/models"
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
