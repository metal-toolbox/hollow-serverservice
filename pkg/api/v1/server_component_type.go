package serverservice

import (
	"go.hollow.sh/serverservice/internal/models"
)

// ServerComponentType provides a way to group server components by the type
type ServerComponentType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (t *ServerComponentType) fromDBModel(dbT *models.ServerComponentType) error {
	t.Name = dbT.Name
	t.ID = dbT.Slug

	return nil
}

func (t *ServerComponentType) toDBModel() (*models.ServerComponentType, error) {
	dbT := &models.ServerComponentType{
		Name: t.Name,
		Slug: t.ID,
	}

	return dbT, nil
}
