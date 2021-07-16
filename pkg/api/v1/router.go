package hollow

import (
	"github.com/gin-gonic/gin"

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
	rg.GET("/hardware/:uuid/versioned-attributes", r.hardwareVersionedAttributesList)
	// rg.PUT("/hardware/:uuid/versioned-attributes", hardwareVersionedAttributesCreate)

	rg.GET("/hardware-component-types", r.hardwareComponentTypeList)
	rg.POST("/hardware-component-types", r.hardwareComponentTypeCreate)
}
