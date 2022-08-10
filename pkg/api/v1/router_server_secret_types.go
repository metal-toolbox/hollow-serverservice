package serverservice

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"go.hollow.sh/serverservice/internal/models"
)

func (r *Router) serverCredentialTypesList(c *gin.Context) {
	pager := parsePagination(c)

	dbTypes, err := models.ServerCredentialTypes(pager.queryMods()...).All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count, err := models.ServerCredentialTypes().Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	types := []ServerCredentialType{}

	for _, dbType := range dbTypes {
		t := ServerCredentialType{}
		t.fromDBModel(dbType)

		types = append(types, t)
	}

	pd := paginationData{
		pageCount:  len(types),
		totalCount: count,
		pager:      pager,
	}

	listResponse(c, types, pd)
}

func (r *Router) serverCredentialTypesCreate(c *gin.Context) {
	var sType models.ServerCredentialType
	if err := c.ShouldBindJSON(&sType); err != nil {
		badRequestResponse(c, "invalid server secret type", err)
		return
	}

	sType.Builtin = false

	if err := sType.Insert(
		c.Request.Context(),
		r.DB,
		boil.Blacklist(models.ServerCredentialTypeColumns.ID),
	); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, sType.Slug)
}
