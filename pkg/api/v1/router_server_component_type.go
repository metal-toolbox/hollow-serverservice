package dcim

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"go.hollow.sh/serverservice/internal/models"
)

func (r *Router) serverComponentTypeCreate(c *gin.Context) {
	var t ServerComponentType
	if err := c.ShouldBindJSON(&t); err != nil {
		badRequestResponse(c, "invalid server component type", err)
		return
	}

	dbT, err := t.toDBModel()
	if err != nil {
		badRequestResponse(c, "invalid server component type", err)
		return
	}

	if err := dbT.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbT.Slug)
}

func (r *Router) serverComponentTypeList(c *gin.Context) {
	pager := parsePagination(c)

	// dbFilter := &gormdb.ServerComponentTypeFilter{
	// 	Name: c.Query("name"),
	// }

	dbTypes, err := models.ServerComponentTypes().All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count, err := models.ServerComponentTypes().Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
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

	pd := paginationData{
		pageCount:  len(types),
		totalCount: count,
		pager:      pager,
	}

	listResponse(c, types, pd)
}
