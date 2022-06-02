package serverservice

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.hollow.sh/toolbox/ginjwt"

	"go.hollow.sh/serverservice/internal/models"
)

// Router provides a router for the v1 API
type Router struct {
	AuthMW *ginjwt.Middleware
	DB     *sqlx.DB
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	amw := r.AuthMW

	rg.GET("/servers", amw.AuthRequired(), amw.RequiredScopes(readScopes("server")), r.serverList)
	rg.POST("/servers", amw.AuthRequired(), amw.RequiredScopes(createScopes("server")), r.serverCreate)
	rg.GET("/servers/:uuid", amw.AuthRequired(), amw.RequiredScopes(readScopes("server")), r.serverGet)
	rg.PUT("/servers/:uuid", amw.AuthRequired(), amw.RequiredScopes(updateScopes("server")), r.serverUpdate)
	rg.DELETE("/servers/:uuid", amw.AuthRequired(), amw.RequiredScopes(deleteScopes("server")), r.serverDelete)

	rg.GET("/servers/:uuid/attributes", amw.AuthRequired(), amw.RequiredScopes(readScopes("server", "server:attributes")), r.serverAttributesList)
	rg.POST("/servers/:uuid/attributes", amw.AuthRequired(), amw.RequiredScopes(createScopes("server", "server:attributes")), r.serverAttributesCreate)
	rg.GET("/servers/:uuid/attributes/:namespace", amw.AuthRequired(), amw.RequiredScopes(readScopes("server", "server:attributes")), r.serverAttributesGet)
	rg.PUT("/servers/:uuid/attributes/:namespace", amw.AuthRequired(), amw.RequiredScopes(updateScopes("server", "server:attributes")), r.serverAttributesUpdate)
	rg.DELETE("/servers/:uuid/attributes/:namespace", amw.AuthRequired(), amw.RequiredScopes(deleteScopes("server", "server:attributes")), r.serverAttributesDelete)

	rg.GET("/servers/:uuid/components", amw.AuthRequired(), amw.RequiredScopes(readScopes("server", "server:component")), r.serverComponentList)
	// rg.POST("/servers/:uuid/components", amw.AuthRequired(), amw.RequiredScopes(createScopes("server", "server:component")))
	// rg.PUT("/servers/:uuid/components", amw.AuthRequired(), amw.RequiredScopes(updateScopes("server", "server:component")))

	rg.GET("/servers/:uuid/versioned-attributes", amw.AuthRequired(), amw.RequiredScopes(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesList)
	rg.POST("/servers/:uuid/versioned-attributes", amw.AuthRequired(), amw.RequiredScopes(createScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesCreate)
	rg.GET("/servers/:uuid/versioned-attributes/:namespace", amw.AuthRequired(), amw.RequiredScopes(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesGet)

	rg.GET("/server-component-types", amw.AuthRequired(), amw.RequiredScopes(readScopes("server-component-types")), r.serverComponentTypeList)
	rg.POST("/server-component-types", amw.AuthRequired(), amw.RequiredScopes(updateScopes("server-component-types")), r.serverComponentTypeCreate)

	rg.GET("/server-component-firmwares", amw.AuthRequired(), amw.RequiredScopes(readScopes("server")), r.serverComponentFirmwareList)
	rg.POST("/server-component-firmwares", amw.AuthRequired(), amw.RequiredScopes(createScopes("server")), r.serverComponentFirmwareCreate)
	rg.GET("/server-component-firmwares/:uuid", amw.AuthRequired(), amw.RequiredScopes(readScopes("server")), r.serverComponentFirmwareGet)
	rg.PUT("/server-component-firmwares/:uuid", amw.AuthRequired(), amw.RequiredScopes(updateScopes("server")), r.serverComponentFirmwareUpdate)
	rg.DELETE("/server-component-firmwares/:uuid", amw.AuthRequired(), amw.RequiredScopes(deleteScopes("server")), r.serverComponentFirmwareDelete)
}

func createScopes(items ...string) []string {
	s := []string{"write", "create"}
	for _, i := range items {
		s = append(s, fmt.Sprintf("create:%s", i))
	}

	return s
}

func readScopes(items ...string) []string {
	s := []string{"read"}
	for _, i := range items {
		s = append(s, fmt.Sprintf("read:%s", i))
	}

	return s
}

func updateScopes(items ...string) []string {
	s := []string{"write", "update"}
	for _, i := range items {
		s = append(s, fmt.Sprintf("update:%s", i))
	}

	return s
}

func deleteScopes(items ...string) []string {
	s := []string{"write", "delete"}
	for _, i := range items {
		s = append(s, fmt.Sprintf("delete:%s", i))
	}

	return s
}

func (r *Router) parseUUID(c *gin.Context) (uuid.UUID, error) {
	u, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		badRequestResponse(c, "failed to parse uuid", err)
	}

	return u, err
}

func (r *Router) loadServerFromParams(c *gin.Context) (*models.Server, error) {
	u, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	srv, err := models.FindServer(c.Request.Context(), r.DB, u.String())
	if err != nil {
		dbErrorResponse(c, err)

		return nil, err
	}

	return srv, nil
}

func (r *Router) loadOrCreateServerFromParams(c *gin.Context) (*models.Server, error) {
	u, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	srv, err := models.FindServer(c.Request.Context(), r.DB, u.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv = &models.Server{ID: u.String()}
			if err := srv.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
				dbErrorResponse(c, err)
				return nil, err
			}

			return srv, nil
		}

		dbErrorResponse(c, err)

		return nil, err
	}

	return srv, nil
}

func (r *Router) loadComponentFirmwareVersionFromParams(c *gin.Context) (*models.ComponentFirmwareVersion, error) {
	u, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	firmware, err := models.FindComponentFirmwareVersion(c.Request.Context(), r.DB, u.String())
	if err != nil {
		dbErrorResponse(c, err)

		return nil, err
	}

	return firmware, nil
}
