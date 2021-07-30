package hollow

import (
	"github.com/gin-gonic/gin"
)

func (r *Router) serverAttributesList(c *gin.Context) {
	u, err := r.parseUUID(c)
	if err != nil {
		return
	}

	pager, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "invalid pagination", err)
		return
	}

	dbAttrs, count, err := r.Store.GetAttributesByServerUUID(u, &pager)
	if err != nil {
		dbFailureResponse(c, err)
		return
	}

	attrs, err := convertFromDBAttributes(dbAttrs)
	if err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	nextCursor := ""

	sz := len(attrs)
	if sz != 0 {
		nextCursor = encodeCursor(attrs[sz-1].CreatedAt)
	}

	pd := paginationData{
		pageCount:  len(attrs),
		totalCount: count,
		nextCursor: nextCursor,
		pager:      pager,
	}

	listResponse(c, attrs, pd)
}

func (r *Router) serverAttributesGet(c *gin.Context) {
	u, err := r.parseUUID(c)
	if err != nil {
		return
	}

	ns := c.Param("namespace")

	dbAttr, err := r.Store.GetAttributesByServerUUIDAndNamespace(u, ns)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	attr := Attributes{}
	if err := attr.fromDBModel(*dbAttr); err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	itemResponse(c, attr)
}

func (r *Router) serverAttributesCreate(c *gin.Context) {
	u, err := r.parseUUID(c)
	if err != nil {
		return
	}

	ns := c.Param("namespace")

	var attr Attributes
	if err := c.ShouldBindJSON(&attr); err != nil {
		badRequestResponse(c, "invalid attributes", err)
		return
	}

	dbAttr, err := attr.toDBModel()
	if err != nil {
		badRequestResponse(c, "invalid attributes", err)
		return
	}

	dbAttr.ServerID = &u
	dbAttr.Namespace = ns

	if err := r.Store.CreateAttributes(&dbAttr); err != nil {
		dbFailureResponse(c, err)
		return
	}

	createdResponse(c, &dbAttr.ID)
}

func (r *Router) serverAttributesUpdate(c *gin.Context) {
	u, err := r.parseUUID(c)
	if err != nil {
		return
	}

	ns := c.Param("namespace")

	var attr Attributes
	if err := c.ShouldBindJSON(&attr); err != nil {
		badRequestResponse(c, "invalid attributes", err)
		return
	}

	err = r.Store.UpdateAttributesByServerUUIDAndNamespace(u, ns, attr.Data)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse(c)
}

func (r *Router) serverAttributesDelete(c *gin.Context) {
	u, err := r.parseUUID(c)
	if err != nil {
		return
	}

	ns := c.Param("namespace")

	dbAttr, err := r.Store.GetAttributesByServerUUIDAndNamespace(u, ns)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	if err = r.Store.DeleteAttributes(dbAttr); err != nil {
		dbFailureResponse(c, err)
		return
	}

	deletedResponse(c)
}
