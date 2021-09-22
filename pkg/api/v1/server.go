package dcim

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.hollow.sh/serverservice/internal/models"
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
	DeletedAt           time.Time             `json:"deleted_at,omitempty"`
}

func (r *Router) getServers(c *gin.Context, params ServerListParams) (models.ServerSlice, int64, error) {
	mods := params.queryMods()

	count, err := models.Servers(mods...).Count(c.Request.Context(), r.DB)
	if err != nil {
		return nil, 0, err
	}

	// add pagination
	mods = append(mods, params.PaginationParams.queryMods()...)

	if params.IncludeDeleted {
		mods = append(mods, qm.WithDeleted())
	}

	s, err := models.Servers(mods...).All(c.Request.Context(), r.DB)
	if err != nil {
		return s, 0, err
	}

	return s, count, nil
}

func (s *Server) fromDBModel(dbS *models.Server) error {
	var err error

	s.UUID, err = uuid.Parse(dbS.ID)
	if err != nil {
		return err
	}

	s.Name = dbS.Name.String
	s.FacilityCode = dbS.FacilityCode.String
	s.CreatedAt = dbS.CreatedAt.Time
	s.UpdatedAt = dbS.UpdatedAt.Time
	s.DeletedAt = dbS.DeletedAt.Time

	if dbS.R != nil {
		if dbS.R.Attributes != nil {
			s.Attributes, err = convertFromDBAttributes(dbS.R.Attributes)
			if err != nil {
				return err
			}
		}

		if dbS.R.ServerComponents != nil {
			s.Components, err = convertDBServerComponents(dbS.R.ServerComponents)
			if err != nil {
				return err
			}
		}

		if dbS.R.VersionedAttributes != nil {
			s.VersionedAttributes, err = convertFromDBVersionedAttributes(dbS.R.VersionedAttributes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) toDBModel() (*models.Server, error) {
	dbS := &models.Server{
		Name:         null.StringFrom(s.Name),
		FacilityCode: null.StringFrom(s.FacilityCode),
	}

	if s.UUID.String() != uuid.Nil.String() {
		dbS.ID = s.UUID.String()
	}

	return dbS, nil
}
