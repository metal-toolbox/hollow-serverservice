package serverservice

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"go.hollow.sh/serverservice/internal/models"
)

func (r *Router) serverSecretTypesList(c *gin.Context) {
	pager := parsePagination(c)

	dbTypes, err := models.ServerSecretTypes(pager.queryMods()...).All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count, err := models.ServerSecretTypes().Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	types := []ServerSecretType{}

	for _, dbType := range dbTypes {
		t := ServerSecretType{}
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

func (r *Router) serverSecretTypesCreate(c *gin.Context) {
	var sType models.ServerSecretType
	if err := c.ShouldBindJSON(&sType); err != nil {
		badRequestResponse(c, "invalid server secret type", err)
		return
	}

	sType.Builtin = false

	if err := sType.Insert(
		c.Request.Context(),
		r.DB,
		boil.Blacklist(models.ServerSecretTypeColumns.ID),
	); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, sType.Slug)
}
