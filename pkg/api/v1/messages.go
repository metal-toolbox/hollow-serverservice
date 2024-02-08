//nolint:wsl,revive
package fleetdbapi

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

var (
	ErrNilServer  = errors.New("bogus server structure provided")
	ErrBadJSONOut = errors.New("object serializaion failed")
	ErrBadJSONIn  = errors.New("object deserializaion failed")
)

// MsgMetadata captures some message-type agnostic descriptive data a consumer might need
type MsgMetadata struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// CreateServer is a message type published via NATS
type CreateServer struct {
	Metadata     *MsgMetadata `json:"metadata,omitempty"`
	Name         null.String  `json:"name"`
	FacilityCode null.String  `json:"facility_code"`
	ID           string       `json:"id"`
}

// NewCreateServerMessage composes a CreateServer message for NATS
func NewCreateServerMessage(srv *models.Server) ([]byte, error) {
	if srv == nil {
		return nil, ErrNilServer
	}
	cs := &CreateServer{
		Metadata: &MsgMetadata{
			CreatedAt: time.Now(),
		},
		Name:         srv.Name,
		FacilityCode: srv.FacilityCode,
		ID:           srv.ID,
	}
	byt, err := json.Marshal(cs)
	if err != nil {
		return nil, errors.Wrap(ErrBadJSONOut, err.Error())
	}
	return byt, err
}

// DeserializeCreateServer reconstitutes a CreateServer from raw bytes
func DeserializeCreateServer(inc []byte) (*CreateServer, error) {
	cs := &CreateServer{}
	if err := json.Unmarshal(inc, cs); err != nil {
		return nil, errors.Wrap(ErrBadJSONIn, err.Error())
	}
	return cs, nil
}
