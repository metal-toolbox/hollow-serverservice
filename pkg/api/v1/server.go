package hollow

import (
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"

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

func (s *Server) fromDBModel(dbS *db.Server) error {
	var err error

	s.UUID, err = uuid.Parse(dbS.ID)
	if err != nil {
		return err
	}

	s.Name = dbS.Name.String
	s.FacilityCode = dbS.FacilityCode.String
	s.CreatedAt = dbS.CreatedAt.Time
	s.UpdatedAt = dbS.UpdatedAt.Time

	s.Attributes, err = convertFromDBAttributes(dbS.R.Attributes)
	if err != nil {
		return err
	}

	s.Components, err = convertDBServerComponents(dbS.R.ServerComponents)
	if err != nil {
		return err
	}

	s.VersionedAttributes, err = convertFromDBVersionedAttributes(dbS.R.VersionedAttributes)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) toDBModel() (*db.Server, error) {
	dbS := &db.Server{
		Name:         null.StringFrom(s.Name),
		FacilityCode: null.StringFrom(s.FacilityCode),
	}

	if s.UUID.String() != uuid.Nil.String() {
		dbS.ID = s.UUID.String()
	}

	// for _, c := range s.Components {
	// 	dbC, err := c.toDBModel(store)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	dbS.ServerComponents = append(dbS.ServerComponents, *dbC)
	// }

	// attrs, err := convertToDBAttributes(s.Attributes)
	// if err != nil {
	// 	return nil, err
	// }

	// dbS.Attributes = attrs

	// verAttrs, err := convertToDBVersionedAttributes(s.VersionedAttributes)
	// if err != nil {
	// 	return nil, err
	// }

	// dbS.VersionedAttributes = verAttrs

	return dbS, nil
}
