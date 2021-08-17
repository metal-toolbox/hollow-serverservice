package hollow

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
	"go.metalkube.net/hollow/internal/gormdb"
	"go.metalkube.net/hollow/pkg/ginjwt"
)

// Router provides a router for the v1 API
type Router struct {
	Store  *gormdb.Store
	AuthMW *ginjwt.Middleware
	DB     *sql.DB
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	amw := r.AuthMW

	rg.GET("/servers", amw.AuthRequired(readScopes("server")), r.serverList)
	rg.POST("/servers", amw.AuthRequired(createScopes("server")), r.serverCreate)
	rg.GET("/servers/:uuid", amw.AuthRequired(readScopes("server")), r.serverGet)
	rg.PUT("/servers/:uuid", amw.AuthRequired(updateScopes("server")), r.serverUpdate)
	rg.DELETE("/servers/:uuid", amw.AuthRequired(deleteScopes("server")), r.serverDelete)

	rg.GET("/servers/:uuid/attributes", amw.AuthRequired(readScopes("server", "server:attributes")), r.serverAttributesList)
	rg.POST("/servers/:uuid/attributes", amw.AuthRequired(createScopes("server", "server:attributes")), r.serverAttributesCreate)
	rg.GET("/servers/:uuid/attributes/:namespace", amw.AuthRequired(readScopes("server", "server:attributes")), r.serverAttributesGet)
	rg.PUT("/servers/:uuid/attributes/:namespace", amw.AuthRequired(updateScopes("server", "server:attributes")), r.serverAttributesUpdate)
	rg.DELETE("/servers/:uuid/attributes/:namespace", amw.AuthRequired(deleteScopes("server", "server:attributes")), r.serverAttributesDelete)

	rg.GET("/servers/:uuid/components", amw.AuthRequired(readScopes("server", "server:component")), r.serverComponentList)
	// rg.POST("/servers/:uuid/components", amw.AuthRequired(createScopes("server", "server:component")))
	// rg.PUT("/servers/:uuid/components", amw.AuthRequired(updateScopes("server", "server:component")))

	rg.GET("/servers/:uuid/versioned-attributes", amw.AuthRequired(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesList)
	rg.POST("/servers/:uuid/versioned-attributes", amw.AuthRequired(createScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesCreate)
	rg.GET("/servers/:uuid/versioned-attributes/:namespace", amw.AuthRequired(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesGet)

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
	u, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	srv, err := db.FindServer(c.Request.Context(), r.DB, u.String())
	if err != nil {
		dbErrorResponse(c, err)
		return nil, err
	}

	return srv, nil
}

func (r *Router) loadOrCreateServerFromParams(c *gin.Context) (*gormdb.Server, error) {
	srvUUID, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	srv, err := r.Store.FindOrCreateServerByUUID(srvUUID)
	if err != nil {
		dbErrorResponse(c, err)
		return nil, err
	}

	return srv, nil
}
