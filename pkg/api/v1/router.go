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
	"gocloud.dev/secrets"

	"go.hollow.sh/serverservice/internal/models"
)

// Router provides a router for the v1 API
type Router struct {
	AuthMW        *ginjwt.Middleware
	DB            *sqlx.DB
	SecretsKeeper *secrets.Keeper
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	amw := r.AuthMW

	// require all calls to have auth
	rg.Use(amw.AuthRequired())

	// /servers
	srvs := rg.Group("/servers")
	{
		srvs.GET("", amw.RequiredScopes(readScopes("server")), r.serverList)
		srvs.POST("", amw.RequiredScopes(createScopes("server")), r.serverCreate)

		srvs.GET("/components", amw.RequiredScopes(readScopes("server:component")), r.serverComponentList)

		// /servers/:uuid
		srv := srvs.Group("/:uuid")
		{
			srv.GET("", amw.RequiredScopes(readScopes("server")), r.serverGet)
			srv.PUT("", amw.RequiredScopes(updateScopes("server")), r.serverUpdate)
			srv.DELETE("", amw.RequiredScopes(deleteScopes("server")), r.serverDelete)

			// /servers/:uuid/attributes
			srvAttrs := srv.Group("/attributes")
			{
				srvAttrs.GET("", amw.RequiredScopes(readScopes("server", "server:attributes")), r.serverAttributesList)
				srvAttrs.POST("", amw.RequiredScopes(createScopes("server", "server:attributes")), r.serverAttributesCreate)
				srvAttrs.GET("/:namespace", amw.RequiredScopes(readScopes("server", "server:attributes")), r.serverAttributesGet)
				srvAttrs.PUT("/:namespace", amw.RequiredScopes(updateScopes("server", "server:attributes")), r.serverAttributesUpdate)
				srvAttrs.DELETE("/:namespace", amw.RequiredScopes(deleteScopes("server", "server:attributes")), r.serverAttributesDelete)
			}

			// /servers/:uuid/components
			srvComponents := srv.Group("/components")
			{
				srvComponents.POST("", amw.RequiredScopes(createScopes("server", "server:component")), r.serverComponentsCreate)
				srvComponents.GET("", amw.RequiredScopes(readScopes("server", "server:component")), r.serverComponentGet)
				srvComponents.PUT("", amw.RequiredScopes(updateScopes("server", "server:component")), r.serverComponentUpdate)
			}

			// /servers/:uuid/secrets/:slug
			svrSecret := srv.Group("secrets/:slug")
			{
				svrSecret.GET("", amw.RequiredScopes([]string{"read:server:secrets"}), r.serverSecretGet)
				svrSecret.PUT("", amw.RequiredScopes([]string{"write:server:secrets"}), r.serverSecretUpsert)
				svrSecret.DELETE("", amw.RequiredScopes([]string{"write:server:secrets"}), r.serverSecretDelete)
			}

			// /servers/:uuid/versioned-attributes
			srvVerAttrs := srv.Group("/versioned-attributes")
			{
				srvVerAttrs.GET("", amw.RequiredScopes(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesList)
				srvVerAttrs.POST("", amw.RequiredScopes(createScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesCreate)
				srvVerAttrs.GET("/:namespace", amw.RequiredScopes(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesGet)
			}
		}
	}

	// /server-component-types
	srvCmpntType := rg.Group("/server-component-types")
	{
		srvCmpntType.GET("", amw.RequiredScopes(readScopes("server-component-types")), r.serverComponentTypeList)
		srvCmpntType.POST("", amw.RequiredScopes(updateScopes("server-component-types")), r.serverComponentTypeCreate)
	}

	// /server-component-firmwares
	srvCmpntFw := rg.Group("/server-component-firmwares")
	{
		srvCmpntFw.GET("", amw.RequiredScopes(readScopes("server")), r.serverComponentFirmwareList)
		srvCmpntFw.POST("", amw.RequiredScopes(createScopes("server")), r.serverComponentFirmwareCreate)
		srvCmpntFw.GET("/:uuid", amw.RequiredScopes(readScopes("server")), r.serverComponentFirmwareGet)
		srvCmpntFw.PUT("/:uuid", amw.RequiredScopes(updateScopes("server")), r.serverComponentFirmwareUpdate)
		srvCmpntFw.DELETE("/:uuid", amw.RequiredScopes(deleteScopes("server")), r.serverComponentFirmwareDelete)
	}

	// /server-secret-types
	srvSecretTypes := rg.Group("/server-secret-types")
	{
		srvSecretTypes.GET("", amw.RequiredScopes(readScopes("server-secret-types")), r.serverSecretTypesList)
		srvSecretTypes.POST("", amw.RequiredScopes(createScopes("server-secret-types")), r.serverSecretTypesCreate)
	}
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
