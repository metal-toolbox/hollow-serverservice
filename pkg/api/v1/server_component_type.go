package hollow

import (
	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// ServerComponentType provides a way to group server components by the type
type ServerComponentType struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

func (t *ServerComponentType) fromDBModel(dbT db.ServerComponentType) error {
	t.UUID = dbT.ID
	t.Name = dbT.Name

	return nil
}

func (t *ServerComponentType) toDBModel() (*db.ServerComponentType, error) {
	dbT := &db.ServerComponentType{
		ID:   t.UUID,
		Name: t.Name,
	}

	return dbT, nil
}
