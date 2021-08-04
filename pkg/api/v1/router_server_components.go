package hollow

import (
	"github.com/gin-gonic/gin"
)

func (r *Router) serverComponentList(c *gin.Context) {
	u, err := r.parseUUID(c)
	if err != nil {
		return
	}

	pager, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "invalid pagination", err)
		return
	}

	dbComps, count, err := r.Store.GetComponentsByServerUUID(u, nil, &pager)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

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
