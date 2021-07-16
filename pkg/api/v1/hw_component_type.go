package hollow

import (
	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// HardwareComponentType provides a way to group hardware components by the type
type HardwareComponentType struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

func (t *HardwareComponentType) fromDBModel(dbT db.HardwareComponentType) error {
	t.UUID = dbT.ID
	t.Name = dbT.Name

	return nil
}

func (t *HardwareComponentType) toDBModel() (*db.HardwareComponentType, error) {
	dbT := &db.HardwareComponentType{
		ID:   t.UUID,
		Name: t.Name,
	}

	return dbT, nil
}
