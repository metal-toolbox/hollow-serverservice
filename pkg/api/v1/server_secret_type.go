package serverservice

import (
	"time"

	"go.hollow.sh/serverservice/internal/models"
)

const (
	// ServerCredentialTypeBMC returns the slug for the builtin ServerCredentialType used
	// to store BMC passwords
	ServerCredentialTypeBMC = "bmc"
)

// ServerCredentialType represents a type of server secret. There are some built in
// default secret types, for example a type exists for BMC passwords.
type ServerCredentialType struct {
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Builtin   bool      `json:"builtin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *ServerCredentialType) fromDBModel(dbT *models.ServerCredentialType) {
	t.Name = dbT.Name
	t.Slug = dbT.Slug
	t.Builtin = dbT.Builtin
	t.CreatedAt = dbT.CreatedAt
	t.UpdatedAt = dbT.UpdatedAt
}
