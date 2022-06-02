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
