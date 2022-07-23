package serverservice

import (
	"time"

	"go.hollow.sh/serverservice/internal/models"
)

const (
	// ServerSecretTypeBMC returns the slug for the builtin ServerSecretType used
	// to store BMC passwords
	ServerSecretTypeBMC = "bmc"
)

// ServerSecretType represents a type of server secret. There are some built in
// default secret types, for example a type exists for BMC passwords.
type ServerSecretType struct {
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Builtin   bool      `json:"builtin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *ServerSecretType) fromDBModel(dbT *models.ServerSecretType) {
	t.Name = dbT.Name
	t.Slug = dbT.Slug
	t.Builtin = dbT.Builtin
	t.CreatedAt = dbT.CreatedAt
	t.UpdatedAt = dbT.UpdatedAt
}
