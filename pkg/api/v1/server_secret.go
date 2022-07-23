package serverservice

import (
	"time"

	"github.com/google/uuid"
)

// ServerSecret provides a way to encrypt secrets about a server in the database
type ServerSecret struct {
	ServerID   uuid.UUID `json:"uuid,omitempty"`
	SecretType string    `json:"secret_type"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type serverSecretValue struct {
	Value string `json:"value"`
}
