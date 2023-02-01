package serverservice

import (
	"go.hollow.sh/serverservice/internal/models"
)

// ServerConditionType provides a way to group server conditions by the type
type ServerConditionType struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

func (t *ServerConditionType) fromDBModel(dbT *models.ServerConditionType) error {
	t.ID = dbT.ID
	t.Slug = dbT.Slug

	return nil
}

func (t *ServerConditionType) toDBModel() (*models.ServerConditionType, error) {
	dbT := &models.ServerConditionType{
		Slug: t.Slug,
	}

	return dbT, nil
}

// ServerConditionTypeSlice is a slice of the ServerConditionType
type ServerConditionTypeSlice []*ServerConditionType

// ByID returns the ServerConditionType matched by its ID field value
func (ts ServerConditionTypeSlice) ByID(id string) *ServerConditionType {
	for _, conditionType := range ts {
		if conditionType.ID == id {
			return conditionType
		}
	}

	return nil
}

// BySlug returns the ServerConditionType matched by its Slug value
func (ts ServerConditionTypeSlice) BySlug(slug string) *ServerConditionType {
	for _, conditionType := range ts {
		if conditionType.Slug == slug {
			return conditionType
		}
	}

	return nil
}
