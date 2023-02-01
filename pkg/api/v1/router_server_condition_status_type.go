package serverservice

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"go.hollow.sh/serverservice/internal/models"
)

func (r *Router) serverConditionStatusTypeCreate(c *gin.Context) {
	var t ServerConditionStatusType
	if err := c.ShouldBindJSON(&t); err != nil {
		badRequestResponse(c, "invalid server condition status type", err)
		return
	}

	dbT, err := t.toDBModel()
	if err != nil {
		badRequestResponse(c, "invalid server condition status type", err)
		return
	}

	if err := dbT.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbT.Slug)
}

func (r *Router) serverConditionStatusTypeList(c *gin.Context) {
	pager := parsePagination(c)

	dbTypes, err := models.ServerConditionStatusTypes().All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count, err := models.ServerConditionStatusTypes().Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	types := []ServerConditionStatusType{}

	for _, dbT := range dbTypes {
		t := ServerConditionStatusType{}
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
