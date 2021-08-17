package hollow

import (
	"github.com/gin-gonic/gin"
)

func (r *Router) serverComponentList(c *gin.Context) {
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	pager, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "invalid pagination", err)
		return
	}

	dbComps, err := srv.ServerComponents().All(c.Request.Context(), r.DB)
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

	nextCursor := ""

	sz := len(comps)
	if sz != 0 {
		nextCursor = encodeCursor(comps[sz-1].CreatedAt)
	}

	pd := paginationData{
		pageCount:  len(comps),
		totalCount: count,
		nextCursor: nextCursor,
		pager:      pager,
	}

	listResponse(c, comps, pd)
}
