package hollow

import (
	"github.com/gin-gonic/gin"
)

// RouteMap will add the routes for this API version to a router group
func RouteMap(rg *gin.RouterGroup) {
	rg.GET("/hardware")
	rg.GET("/hardware/:uuid/bios-configs", hardwareBIOSConfigList)
	rg.POST("/bios-config", biosConfigCreate)
}
