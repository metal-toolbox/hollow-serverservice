package serverservice

import (
	"go.hollow.sh/serverservice/internal/models"
)

// ServerConditionStatusType provides a way to group server conditions statuses by the type
type ServerConditionStatusType struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

func (t *ServerConditionStatusType) fromDBModel(dbT *models.ServerConditionStatusType) error {
	t.ID = dbT.ID
	t.Slug = dbT.Slug

	return nil
}

func (t *ServerConditionStatusType) toDBModel() (*models.ServerConditionStatusType, error) {
	dbT := &models.ServerConditionStatusType{
		Slug: t.Slug,
	}

	return dbT, nil
}

// ServerConditionStatusTypeSlice is a slice of the ServerConditionStatusType
type ServerConditionStatusTypeSlice []*ServerConditionStatusType

// ByID returns the ServerConditionType matched by its ID field value
func (ts ServerConditionStatusTypeSlice) ByID(id string) *ServerConditionStatusType {
	for _, conditionStatusType := range ts {
		if conditionStatusType.ID == id {
			return conditionStatusType
		}
	}

	return nil
}

// BySlug returns the ServerConditionStatusType matched by its Slug value
func (ts ServerConditionStatusTypeSlice) BySlug(slug string) *ServerConditionStatusType {
	for _, conditionStatusType := range ts {
		if conditionStatusType.Slug == slug {
			return conditionStatusType
		}
	}

	return nil
}
