package serverservice

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (r *Router) serverComponentList(c *gin.Context) {
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	pager := parsePagination(c)

	dbComps, err := srv.ServerComponents(qm.Load("ServerComponentType")).All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count := int64(0)

	comps, err := convertDBServerComponents(dbComps)
	if err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	pd := paginationData{
		pageCount:  len(comps),
		totalCount: count,
		pager:      pager,
	}

	listResponse(c, comps, pd)
}
