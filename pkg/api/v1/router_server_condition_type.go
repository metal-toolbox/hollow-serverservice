package serverservice

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"go.hollow.sh/serverservice/internal/models"
)

func (r *Router) serverConditionTypeCreate(c *gin.Context) {
	var t ServerConditionType
	if err := c.ShouldBindJSON(&t); err != nil {
		badRequestResponse(c, "invalid server condition type", err)
		return
	}

	dbT, err := t.toDBModel()
	if err != nil {
		badRequestResponse(c, "invalid server condition type", err)
		return
	}

	if err := dbT.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbT.Slug)
}

func (r *Router) serverConditionTypeList(c *gin.Context) {
	pager := parsePagination(c)

	dbTypes, err := models.ServerConditionTypes().All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count, err := models.ServerConditionTypes().Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	types := []ServerConditionType{}

	for _, dbT := range dbTypes {
		t := ServerConditionType{}
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
