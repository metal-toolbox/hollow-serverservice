package hollow

import (
	"errors"
	"net/http"
	"strconv"

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

	rg.GET("/servers", amw.AuthRequired([]string{"read"}), r.serverList)
	rg.POST("/servers", amw.AuthRequired([]string{"create", "write"}), r.serverCreate)
	rg.GET("/servers/:uuid", amw.AuthRequired([]string{"read"}), r.serverGet)
	rg.DELETE("/servers/:uuid", amw.AuthRequired([]string{"write"}), r.serverDelete)
	rg.GET("/servers/:uuid/versioned-attributes", amw.AuthRequired([]string{"read"}), r.serverVersionedAttributesList)
	rg.PUT("/servers/:uuid/versioned-attributes", amw.AuthRequired([]string{"create", "create:versionedattributes", "write"}), r.serverVersionedAttributesCreate)

	rg.GET("/server-component-types", amw.AuthRequired([]string{"read"}), r.serverComponentTypeList)
	rg.POST("/server-component-types", amw.AuthRequired([]string{"create", "write"}), r.serverComponentTypeCreate)
}

func parsePagination(c *gin.Context) db.Pagination {
	// Initializing default
	limit := db.DefaultPaginationSize
	page := 1
	sort := ""
	query := c.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		case "sort":
			sort = queryValue
		}
	}

	return db.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

func (r *Router) loadServerFromParams(c *gin.Context) (*db.Server, error) {
	srvUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid server uuid", "error": err.Error()})
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
	srvUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid server uuid", "error": err.Error()})
		return nil, err
	}

	srv, err := r.Store.FindOrCreateServerByUUID(srvUUID)
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
