package hollow

import (
	"go.metalkube.net/hollow/internal/db"
)

// ServerComponentType provides a way to group server components by the type
type ServerComponentType struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func (t *ServerComponentType) fromDBModel(dbT db.ServerComponentType) error {
	t.Name = dbT.Name
	t.Slug = dbT.Slug

	return nil
}

func (t *ServerComponentType) toDBModel() (*db.ServerComponentType, error) {
	dbT := &db.ServerComponentType{
		Name: t.Name,
		Slug: t.Slug,
	}

	return dbT, nil
}
