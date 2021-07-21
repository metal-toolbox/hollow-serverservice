package hollow

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go.metalkube.net/hollow/internal/db"
)

func (r *Router) serverComponentTypeCreate(c *gin.Context) {
	var t ServerComponentType
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid server component type",
			"error":   err.Error(),
		})

		return
	}

	dbT, err := t.toDBModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid server component type", "error": err.Error()})
		return
	}

	if err := r.Store.CreateServerComponentType(dbT); err != nil {
		dbFailureResponse(c, err)
		return
	}

	createdResponse(c, &dbT.ID)
}

func (r *Router) serverComponentTypeList(c *gin.Context) {
	pager := parsePagination(c)

	dbFilter := &db.ServerComponentTypeFilter{
		Name: c.Query("name"),
	}

	dbTypes, err := r.Store.GetServerComponentTypes(dbFilter, &pager)
	if err != nil {
		dbFailureResponse(c, err)
		return
	}

	types := []ServerComponentType{}

	for _, dbT := range dbTypes {
		t := ServerComponentType{}
		if err := t.fromDBModel(dbT); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		types = append(types, t)
	}

	c.JSON(http.StatusOK, types)
}
