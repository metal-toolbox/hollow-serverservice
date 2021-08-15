package hollow

import "go.metalkube.net/hollow/internal/gormdb"

// ServerComponentType provides a way to group server components by the type
type ServerComponentType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (t *ServerComponentType) fromDBModel(dbT gormdb.ServerComponentType) error {
	t.Name = dbT.Name
	t.ID = dbT.Slug

	return nil
}

func (t *ServerComponentType) toDBModel() (*gormdb.ServerComponentType, error) {
	dbT := &gormdb.ServerComponentType{
		Name: t.Name,
		Slug: t.ID,
	}

	return dbT, nil
}
