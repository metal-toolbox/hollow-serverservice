package serverservice

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.hollow.sh/toolbox/events"
	"go.hollow.sh/toolbox/ginauth"
	"go.uber.org/zap"
	"gocloud.dev/secrets"

	"go.hollow.sh/serverservice/internal/models"
)

// Router provides a router for the v1 API
type Router struct {
	AuthMW        *ginauth.MultiTokenMiddleware
	DB            *sqlx.DB
	SecretsKeeper *secrets.Keeper
	Logger        *zap.Logger
	EventStream   events.Stream
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	amw := r.AuthMW

	// /servers
	srvs := rg.Group("/servers")
	{
		srvs.GET("", amw.AuthRequired(readScopes("server")), r.serverList)
		srvs.POST("", amw.AuthRequired(createScopes("server")), r.serverCreate)

		srvs.GET("/components", amw.AuthRequired(readScopes("server:component")), r.serverComponentList)

		// /servers/:uuid
		srv := srvs.Group("/:uuid")
		{
			srv.GET("", amw.AuthRequired(readScopes("server")), r.serverGet)
			srv.PUT("", amw.AuthRequired(updateScopes("server")), r.serverUpdate)
			srv.DELETE("", amw.AuthRequired(deleteScopes("server")), r.serverDelete)

			// /servers/:uuid/attributes
			srvAttrs := srv.Group("/attributes")
			{
				srvAttrs.GET("", amw.AuthRequired(readScopes("server", "server:attributes")), r.serverAttributesList)
				srvAttrs.POST("", amw.AuthRequired(createScopes("server", "server:attributes")), r.serverAttributesCreate)
				srvAttrs.GET("/:namespace", amw.AuthRequired(readScopes("server", "server:attributes")), r.serverAttributesGet)
				srvAttrs.PUT("/:namespace", amw.AuthRequired(updateScopes("server", "server:attributes")), r.serverAttributesUpdate)
				srvAttrs.DELETE("/:namespace", amw.AuthRequired(deleteScopes("server", "server:attributes")), r.serverAttributesDelete)
			}

			// /servers/:uuid/components
			srvComponents := srv.Group("/components")
			{
				srvComponents.POST("", amw.AuthRequired(createScopes("server", "server:component")), r.serverComponentsCreate)
				srvComponents.GET("", amw.AuthRequired(readScopes("server", "server:component")), r.serverComponentGet)
				srvComponents.PUT("", amw.AuthRequired(updateScopes("server", "server:component")), r.serverComponentUpdate)
				srvComponents.DELETE("", amw.AuthRequired(deleteScopes("server", "server:component")), r.serverComponentDelete)
			}

			// /servers/:uuid/credentials/:slug
			svrCreds := srv.Group("credentials/:slug")
			{
				svrCreds.GET("", amw.AuthRequired([]string{"read:server:credentials"}), r.serverCredentialGet)
				svrCreds.PUT("", amw.AuthRequired([]string{"write:server:credentials"}), r.serverCredentialUpsert)
				svrCreds.DELETE("", amw.AuthRequired([]string{"write:server:credentials"}), r.serverCredentialDelete)
			}

			// /servers/:uuid/versioned-attributes
			srvVerAttrs := srv.Group("/versioned-attributes")
			{
				srvVerAttrs.GET("", amw.AuthRequired(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesList)
				srvVerAttrs.POST("", amw.AuthRequired(createScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesCreate)
				srvVerAttrs.GET("/:namespace", amw.AuthRequired(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesGet)
			}
		}
	}

	// /server-component-types
	srvCmpntType := rg.Group("/server-component-types")
	{
		srvCmpntType.GET("", amw.AuthRequired(readScopes("server-component-types")), r.serverComponentTypeList)
		srvCmpntType.POST("", amw.AuthRequired(updateScopes("server-component-types")), r.serverComponentTypeCreate)
	}

	// /server-component-firmwares
	srvCmpntFw := rg.Group("/server-component-firmwares")
	{
		srvCmpntFw.GET("", amw.AuthRequired(readScopes("server-component-firmwares")), r.serverComponentFirmwareList)
		srvCmpntFw.POST("", amw.AuthRequired(createScopes("server-component-firmwares")), r.serverComponentFirmwareCreate)
		srvCmpntFw.GET("/:uuid", amw.AuthRequired(readScopes("server-component-firmwares")), r.serverComponentFirmwareGet)
		srvCmpntFw.PUT("/:uuid", amw.AuthRequired(updateScopes("server-component-firmwares")), r.serverComponentFirmwareUpdate)
		srvCmpntFw.DELETE("/:uuid", amw.AuthRequired(deleteScopes("server-component-firmwares")), r.serverComponentFirmwareDelete)
	}

	// /server-credential-types
	srvCredentialTypes := rg.Group("/server-credential-types")
	{
		srvCredentialTypes.GET("", amw.AuthRequired(readScopes("server-credential-types")), r.serverCredentialTypesList)
		srvCredentialTypes.POST("", amw.AuthRequired(createScopes("server-credential-types")), r.serverCredentialTypesCreate)
	}

	// /server-component-firmware-sets
	srvCmpntFwSets := rg.Group("/server-component-firmware-sets")
	{
		srvCmpntFwSets.GET("", amw.AuthRequired(readScopes("server-component-firmware-sets")), r.serverComponentFirmwareSetList)
		srvCmpntFwSets.POST("", amw.AuthRequired(createScopes("server-component-firmware-sets")), r.serverComponentFirmwareSetCreate)
		srvCmpntFwSets.GET("/:uuid", amw.AuthRequired(readScopes("server-component-firmware-sets")), r.serverComponentFirmwareSetGet)
		srvCmpntFwSets.PUT("/:uuid", amw.AuthRequired(updateScopes("server-component-firmware-sets")), r.serverComponentFirmwareSetUpdate)
		srvCmpntFwSets.DELETE("/:uuid", amw.AuthRequired(deleteScopes("server-component-firmware-sets")), r.serverComponentFirmwareSetDelete)
		srvCmpntFwSets.POST("/:uuid/remove-firmware", amw.AuthRequired(deleteScopes("server-component-firmware-sets")), r.serverComponentFirmwareSetRemoveFirmware)
	}

	// /bill-of-materials
	srvBoms := rg.Group("/bill-of-materials")
	{
		// /bill-of-materials/batch-boms-upload
		uploadFile := srvBoms.Group("/batch-upload")
		{
			uploadFile.POST("", amw.AuthRequired(createScopes("batch-upload")), r.bomsUpload)
		}

		// /bill-of-materials/aoc-mac-address
		srvBomByAocMacAddress := srvBoms.Group("/aoc-mac-address")
		{
			srvBomByAocMacAddress.GET("/:aoc_mac_address", amw.AuthRequired(readScopes("aoc-mac-address")), r.getBomFromAocMacAddress)
		}

		// /bill-of-materials/bmc-mac-address
		srvBomByBmcMacAddress := srvBoms.Group("/bmc-mac-address")
		{
			srvBomByBmcMacAddress.GET("/:bmc_mac_address", amw.AuthRequired(readScopes("bmc-mac-address")), r.getBomFromBmcMacAddress)
		}
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
		return nil, errors.Wrap(ErrUUIDParse, err.Error())
	}

	srv, err := models.FindServer(c.Request.Context(), r.DB, u.String())
	if err != nil {
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

// publish a CreateServer message to the event stream. if the publish fails...?
//
//nolint:wsl
func (r *Router) publishCreateServerMessage(ctx context.Context, srv *models.Server) {
	if r.EventStream == nil {
		r.Logger.Error("Event publish skipped, eventStream not connected")
		return
	}
	subject := strings.Join([]string{"server", "create"}, ".")
	payload, err := NewCreateServerMessage(srv)
	if err != nil {
		r.Logger.With(zap.Error(err)).Error("unable to create a create-server message")
		return
	}
	if err := r.EventStream.Publish(ctx, subject, payload); err != nil {
		r.Logger.With(zap.Error(err)).Error("unable to publish create-server message")
		return
	}
}
