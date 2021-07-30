package hollow

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
	"go.metalkube.net/hollow/pkg/ginjwt"
)

// Router provides a router for the v1 API
type Router struct {
	Store  *db.Store
	AuthMW *ginjwt.Middleware
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	amw := r.AuthMW

	rg.GET("/servers", amw.AuthRequired(readScopes("servers")), r.serverList)
	rg.POST("/servers", amw.AuthRequired(createScopes("servers")), r.serverCreate)
	rg.GET("/servers/:uuid", amw.AuthRequired(readScopes("servers")), r.serverGet)
	rg.DELETE("/servers/:uuid", amw.AuthRequired(deleteScopes("servers")), r.serverDelete)

	rg.GET("/servers/:uuid/attributes/", amw.AuthRequired(readScopes("servers", "servers:attributes")), r.serverAttributesList)
	rg.POST("/servers/:uuid/attributes/", amw.AuthRequired(createScopes("servers", "servers:attributes")), r.serverAttributesCreate)
	rg.GET("/servers/:uuid/attributes/:namespace", amw.AuthRequired(readScopes("servers", "servers:attributes")), r.serverAttributesGet)
	rg.PUT("/servers/:uuid/attributes/:namespace", amw.AuthRequired(updateScopes("servers", "servers:attributes")), r.serverAttributesUpdate)
	rg.DELETE("/servers/:uuid/attributes/:namespace", amw.AuthRequired(deleteScopes("servers", "servers:attributes")), r.serverAttributesDelete)

	rg.GET("/servers/:uuid/versioned-attributes", amw.AuthRequired(readScopes("servers", "servers:versioned-attributes")), r.serverVersionedAttributesList)
	rg.POST("/servers/:uuid/versioned-attributes", amw.AuthRequired(createScopes("servers", "servers:versioned-attributes")), r.serverVersionedAttributesCreate)

	rg.GET("/server-component-types", amw.AuthRequired(readScopes("server-component-types")), r.serverComponentTypeList)
	rg.POST("/server-component-types", amw.AuthRequired(updateScopes("server-component-types")), r.serverComponentTypeCreate)
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

func (r *Router) loadServerFromParams(c *gin.Context) (*db.Server, error) {
	srvUUID, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	srv, err := r.Store.FindServerByUUID(srvUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			notFoundResponse(c, err)
			return nil, err
		}

		dbFailureResponse(c, err)

		return nil, err
	}

	return srv, nil
}

func (r *Router) loadOrCreateServerFromParams(c *gin.Context) (*db.Server, error) {
	srvUUID, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	srv, err := r.Store.FindOrCreateServerByUUID(srvUUID)
	if err != nil {
		dbFailureResponse(c, err)
		return nil, err
	}

	return srv, nil
}
