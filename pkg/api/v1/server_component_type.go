package serverservice

import (
	"go.hollow.sh/serverservice/internal/models"
)

// ServerComponentType provides a way to group server components by the type
type ServerComponentType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (t *ServerComponentType) fromDBModel(dbT *models.ServerComponentType) error {
	t.ID = dbT.ID
	t.Name = dbT.Name
	t.Slug = dbT.Slug

	return nil
}

func (t *ServerComponentType) toDBModel() (*models.ServerComponentType, error) {
	dbT := &models.ServerComponentType{
		ID:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}

	return dbT, nil
}

// ServerComponentTypeSlice is a slice of the ServerComponentType
type ServerComponentTypeSlice []*ServerComponentType

// ByID returns the ServerComponentType matched by its ID field value
func (ts ServerComponentTypeSlice) ByID(id string) *ServerComponentType {
	for _, componentType := range ts {
		if componentType.ID == id {
			return componentType
		}
	}

	return nil
}

// ByName returns the ServerComponentType matched by its Name field value
func (ts ServerComponentTypeSlice) ByName(name string) *ServerComponentType {
	for _, componentType := range ts {
		if componentType.Name == name {
			return componentType
		}
	}

	return nil
}

// BySlug returns the ServerComponentType matched by its Slug field value
func (ts ServerComponentTypeSlice) BySlug(slug string) *ServerComponentType {
	for _, componentType := range ts {
		if componentType.Slug == slug {
			return componentType
		}
	}

	return nil
}
