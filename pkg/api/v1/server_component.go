package serverservice

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerComponent represents a component of a server. These can be things like
// processors, NICs, hard drives, etc.
//
// Note: when setting validator struct tags, ensure no extra spaces are present between
//
//	comma separated values or validation will fail with a not so useful 500 error.
type ServerComponent struct {
	UUID                uuid.UUID             `json:"uuid"`
	ServerUUID          uuid.UUID             `json:"server_uuid" binding:"required"`
	Name                string                `json:"name" binding:"required"`
	Vendor              string                `json:"vendor"`
	Model               string                `json:"model"`
	Serial              string                `json:"serial" binding:"required"`
	Attributes          []Attributes          `json:"attributes"`
	VersionedAttributes []VersionedAttributes `json:"versioned_attributes"`
	ComponentTypeID     string                `json:"component_type_id" binding:"required"`
	ComponentTypeName   string                `json:"component_type_name" binding:"required"`
	ComponentTypeSlug   string                `json:"component_type_slug" binding:"required"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}

// ServerComponentSlice is a slice of ServerComponent objects
type ServerComponentSlice []ServerComponent

func convertDBServerComponents(dbComponents models.ServerComponentSlice) ([]ServerComponent, error) {
	components := []ServerComponent{}
	if dbComponents == nil {
		return components, nil
	}

	for _, dbC := range dbComponents {
		c := ServerComponent{}
		if err := c.fromDBModel(dbC); err != nil {
			return nil, err
		}

		components = append(components, c)
	}

	return components, nil
}

// getServerComponents returns server components based on query parameters
func (r *Router) getServerComponents(c *gin.Context, params []ServerComponentListParams, pagination PaginationParams) (models.ServerComponentSlice, int64, error) {
	mods := []qm.QueryMod{}
	// TODO(joel): is there a table name const we could use?
	tableName := "server_components"

	// for each parameter, setup the query modifiers
	for _, param := range params {
		mods = append(mods, param.queryMods(tableName))
	}

	count, err := models.ServerComponents(mods...).Count(c.Request.Context(), r.DB)
	if err != nil {
		return nil, 0, err
	}

	// add pagination
	mods = append(mods, pagination.serverComponentsQueryMods()...)

	sc, err := models.ServerComponents(mods...).All(c.Request.Context(), r.DB)
	if err != nil {
		return sc, 0, err
	}

	return sc, count, nil
}

// fromDBModel populates the ServerComponent object fields based on values in the store
func (c *ServerComponent) fromDBModel(dbC *models.ServerComponent) error {
	var err error

	c.UUID, err = uuid.Parse(dbC.ID)
	if err != nil {
		return err
	}

	c.ServerUUID, err = uuid.Parse(dbC.ServerID)
	if err != nil {
		return err
	}

	c.Name = dbC.Name.String
	c.Vendor = dbC.Vendor.String
	c.Model = dbC.Model.String
	c.Serial = dbC.Serial.String
	c.CreatedAt = dbC.CreatedAt.Time
	c.UpdatedAt = dbC.UpdatedAt.Time

	if dbC.R != nil && dbC.R.ServerComponentType != nil {
		c.ComponentTypeID = dbC.R.ServerComponentType.ID
		c.ComponentTypeName = dbC.R.ServerComponentType.Name
		c.ComponentTypeSlug = dbC.R.ServerComponentType.Slug
	}

	// relation attributes
	if dbC.R.Attributes != nil {
		c.Attributes, err = convertFromDBAttributes(dbC.R.Attributes)
		if err != nil {
			return err
		}
	}

	// relation versioned attributes
	if dbC.R.VersionedAttributes != nil {
		c.VersionedAttributes, err = convertFromDBVersionedAttributes(dbC.R.VersionedAttributes)
		if err != nil {
			return err
		}
	}

	return nil
}

// toDBModel converts a ServerComponent object to a model.ServerComponent object
func (c *ServerComponent) toDBModel(serverID string) *models.ServerComponent {
	return &models.ServerComponent{
		ID:                    c.UUID.String(),
		ServerID:              serverID,
		ServerComponentTypeID: c.ComponentTypeID,
		Name:                  null.StringFrom(c.Name),
		Vendor:                null.StringFrom(c.Vendor),
		Model:                 null.StringFrom(c.Model),
		Serial:                null.StringFrom(c.Serial),
	}
}
