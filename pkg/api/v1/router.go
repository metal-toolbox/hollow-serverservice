package hollow

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// Router provides a router for the v1 API
type Router struct {
	Store *db.Store
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	rg.GET("/hardware", r.hardwareList)
	rg.POST("/hardware", r.hardwareCreate)
	rg.GET("/hardware/:uuid", r.hardwareGet)
	rg.DELETE("/hardware/:uuid", r.hardwareDelete)
	rg.GET("/hardware/:uuid/versioned-attributes", r.hardwareVersionedAttributesList)
	rg.PUT("/hardware/:uuid/versioned-attributes", r.hardwareVersionedAttributesCreate)

	rg.GET("/hardware-component-types", r.hardwareComponentTypeList)
	rg.POST("/hardware-component-types", r.hardwareComponentTypeCreate)
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
