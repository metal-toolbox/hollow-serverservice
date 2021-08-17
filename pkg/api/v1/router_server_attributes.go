package hollow

import (
	"github.com/gin-gonic/gin"

	"go.metalkube.net/hollow/internal/db"
)

func (r *Router) serverAttributesList(c *gin.Context) {
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	pager, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "invalid pagination", err)
		return
	}

	dbAttrs, err := srv.Attributes().All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count := int64(0)

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
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	ns := c.Param("namespace")

	dbAttr, err := srv.Attributes(db.AttributeWhere.Namespace.EQ(ns)).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	attr := Attributes{}
	if err := attr.fromDBModel(dbAttr); err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	itemResponse(c, attr)
}

func (r *Router) serverAttributesCreate(c *gin.Context) {
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

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

	if err := srv.AddAttributes(c.Request.Context(), r.DB, true, dbAttr); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbAttr.Namespace)
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

	updatedResponse(c, ns)
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
		dbErrorResponse(c, err)
		return
	}

	deletedResponse(c)
}
