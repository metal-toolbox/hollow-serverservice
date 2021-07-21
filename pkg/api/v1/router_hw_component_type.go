package hollow

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go.metalkube.net/hollow/internal/db"
)

func (r *Router) hardwareComponentTypeCreate(c *gin.Context) {
	var t HardwareComponentType
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid hardware component type",
			"error":   err.Error(),
		})

		return
	}

	dbT, err := t.toDBModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid hardware component type", "error": err.Error()})
		return
	}

	if err := r.Store.CreateHardwareComponentType(dbT); err != nil {
		dbFailureResponse(c, err)
		return
	}

	createdResponse(c, &dbT.ID)
}

func (r *Router) hardwareComponentTypeList(c *gin.Context) {
	pager := parsePagination(c)

	dbFilter := &db.HardwareComponentTypeFilter{
		Name: c.Query("name"),
	}

	dbTypes, err := r.Store.GetHardwareComponentTypes(dbFilter, &pager)
	if err != nil {
		dbFailureResponse(c, err)
		return
	}

	types := []HardwareComponentType{}

	for _, dbT := range dbTypes {
		t := HardwareComponentType{}
		if err := t.fromDBModel(dbT); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		types = append(types, t)
	}

	c.JSON(http.StatusOK, types)
}
