package hollow

import (
	"time"

	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// Server represents a server in a facility
type Server struct {
	UUID                uuid.UUID             `json:"uuid"`
	Name                string                `json:"name"`
	FacilityCode        string                `json:"facility"`
	Attributes          []Attributes          `json:"attributes"`
	Components          []ServerComponent     `json:"components"`
	VersionedAttributes []VersionedAttributes `json:"versioned_attributes"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}

func (s *Server) fromDBModel(dbS db.Server) error {
	var err error

	s.UUID = dbS.ID
	s.Name = dbS.Name
	s.FacilityCode = dbS.FacilityCode
	s.CreatedAt = dbS.CreatedAt
	s.UpdatedAt = dbS.UpdatedAt

	s.Attributes, err = convertFromDBAttributes(dbS.Attributes)
	if err != nil {
		return err
	}

	s.Components, err = convertDBServerComponents(dbS.ServerComponents)
	if err != nil {
		return err
	}

	s.VersionedAttributes, err = convertFromDBVersionedAttributes(dbS.VersionedAttributes)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) toDBModel() (*db.Server, error) {
	dbS := &db.Server{
		ID:           s.UUID,
		Name:         s.Name,
		FacilityCode: s.FacilityCode,
	}

	for _, c := range s.Components {
		dbC, err := c.toDBModel()
		if err != nil {
			return nil, err
		}

		dbS.ServerComponents = append(dbS.ServerComponents, *dbC)
	}

	attrs, err := convertToDBAttributes(s.Attributes)
	if err != nil {
		return nil, err
	}

	dbS.Attributes = attrs

	verAttrs, err := convertToDBVersionedAttributes(s.VersionedAttributes)
	if err != nil {
		return nil, err
	}

	dbS.VersionedAttributes = verAttrs

	return dbS, nil
}
