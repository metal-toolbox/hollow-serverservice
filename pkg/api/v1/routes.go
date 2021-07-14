package hollow

import (
	"github.com/gin-gonic/gin"
)

// RouteMap will add the routes for this API version to a router group
func RouteMap(rg *gin.RouterGroup) {
	rg.POST("/bios-config", biosConfigCreate)

	rg.GET("/hardware")
	rg.POST("/hardware", hardwareCreate)
	rg.GET("/hardware/:uuid", hardwareGet)
	rg.GET("/hardware/:uuid/bios-configs", hardwareBIOSConfigList)

	rg.GET("/hardware-component-types", hardwareComponentTypeList)
	rg.POST("/hardware-component-types", hardwareComponentTypeCreate)
}
