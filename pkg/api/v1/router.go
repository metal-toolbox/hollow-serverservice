package hollow

import (
	"errors"
	"net/http"

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

	rg.GET("/hardware", amw.AuthRequired([]string{"read"}), r.hardwareList)
	rg.POST("/hardware", amw.AuthRequired([]string{"create", "write"}), r.hardwareCreate)
	rg.GET("/hardware/:uuid", amw.AuthRequired([]string{"read"}), r.hardwareGet)
	rg.DELETE("/hardware/:uuid", amw.AuthRequired([]string{"write"}), r.hardwareDelete)
	rg.GET("/hardware/:uuid/versioned-attributes", amw.AuthRequired([]string{"read"}), r.hardwareVersionedAttributesList)
	rg.PUT("/hardware/:uuid/versioned-attributes", amw.AuthRequired([]string{"create", "create:versionedattributes", "write"}), r.hardwareVersionedAttributesCreate)

	rg.GET("/hardware-component-types", amw.AuthRequired([]string{"read"}), r.hardwareComponentTypeList)
	rg.POST("/hardware-component-types", amw.AuthRequired([]string{"create", "write"}), r.hardwareComponentTypeCreate)
}

func (r *Router) loadHardwareFromParams(c *gin.Context) (*db.Hardware, error) {
	hwUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid hardware uuid", "error": err.Error()})
		return nil, err
	}

	hw, err := r.Store.GetHardwareByUUID(hwUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			notFoundResponse(c, err)
			return nil, err
		}

		dbFailureResponse(c, err)

		return nil, err
	}

	return hw, nil
}
