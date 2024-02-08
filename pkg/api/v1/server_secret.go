package fleetdbapi

import (
	"time"

	"github.com/google/uuid"
)

// ServerCredential provides a way to encrypt secrets about a server in the database
type ServerCredential struct {
	ServerID   uuid.UUID `json:"uuid,omitempty"`
	SecretType string    `json:"secret_type"`
	Password   string    `json:"password"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type serverCredentialValues struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
