package serverservice

import (
	"encoding/json"
	"time"
)

// ServerCondition holds information of a condition set on a server.
//
// This condition in turn is expected to be acted on by a controller external to serverservice.
type ServerCondition struct {
	Slug         string          `json:"slug"`
	Status       string          `json:"status"`
	Parameters   json.RawMessage `json:"parameters"`
	StatusOutput json.RawMessage `json:"status_output"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ServerConditionSlice is a slice of Conditions associated with a server.
type ServerConditionSlice []*ServerCondition
