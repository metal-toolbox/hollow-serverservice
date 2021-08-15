package hollow

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "github.com/volatiletech/null/v8"
	// "github.com/volatiletech/sqlboiler/v4/boil"
	// "go.metalkube.net/hollow/models"
)

func (r *Router) serverList(c *gin.Context) {
	pager, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "invalid pagination", err)
		return
	}

	var params ServerListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		badRequestResponse(c, "invalid filter", err)
		return
	}

	alp, err := parseQueryAttributesListParams(c, "attr")
	if err != nil {
		badRequestResponse(c, "invalid attributes list params", err)
		return
	}

	params.AttributeListParams = alp

	valp, err := parseQueryAttributesListParams(c, "ver_attr")
	if err != nil {
		badRequestResponse(c, "invalid versioned attributes list params", err)
		return
	}

	params.VersionedAttributeListParams = valp

	sclp, err := parseQueryServerComponentsListParams(c)
	if err != nil {
		badRequestResponse(c, "invalid server component list params", err)
		return
	}

	params.ComponentListParams = sclp

	dbFilter, err := params.dbFilter(r)
	if err != nil {
		badRequestResponse(c, "invalid list params", err)
		return
	}

	dbSRV, count, err := r.Store.GetServers(dbFilter, &pager)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	srvs := []Server{}

	for _, dbS := range dbSRV {
		s := Server{}
		if err := s.fromDBModel(dbS); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		srvs = append(srvs, s)
	}

	nextCursor := ""

	sz := len(srvs)
	if sz != 0 {
		nextCursor = encodeCursor(srvs[sz-1].CreatedAt)
	}

	pd := paginationData{
		pageCount:  len(srvs),
		totalCount: count,
		nextCursor: nextCursor,
		pager:      pager,
	}

	listResponse(c, srvs, pd)
}

func (r *Router) serverGet(c *gin.Context) {
	dbSRV, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	var srv Server
	if err = srv.fromDBModel(*dbSRV); err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	itemResponse(c, srv)
}

func (r *Router) serverCreate(c *gin.Context) {
	var srv Server
	if err := c.ShouldBindJSON(&srv); err != nil {
		badRequestResponse(c, "invalid server", err)
		return
	}

	dbSRV, err := srv.toDBModel(r.Store)
	if err != nil {
		badRequestResponse(c, "invalid server", err)
		return
	}

	if err := r.Store.CreateServer(dbSRV); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbSRV.ID.String())
}

func (r *Router) serverDelete(c *gin.Context) {
	dbSRV, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	if err = r.Store.DeleteServer(dbSRV); err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse(c)
}

func (r *Router) serverUpdate(c *gin.Context) {
	u, err := r.parseUUID(c)
	if err != nil {
		return
	}

	var srv Server
	if err := c.ShouldBindJSON(&srv); err != nil {
		badRequestResponse(c, "invalid server", err)
		return
	}

	dbSRV, err := srv.toDBModel(r.Store)
	if err != nil {
		badRequestResponse(c, "invalid server", err)
		return
	}

	if err := r.Store.UpdateServer(u, *dbSRV); err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, u.String())
}

func (r *Router) serverVersionedAttributesGet(c *gin.Context) {
	pager, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "invalid pagination", err)
		return
	}

	srvUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		badRequestResponse(c, "invalid server uuid", err)
		return
	}

	ns := c.Param("namespace")

	dbVA, count, err := r.Store.GetVersionedAttributes(srvUUID, ns, &pager)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	va := []VersionedAttributes{}

	for _, dbA := range dbVA {
		a := VersionedAttributes{}
		if err := a.fromDBModel(dbA); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		va = append(va, a)
	}

	nextCursor := ""

	sz := len(va)
	if sz != 0 {
		nextCursor = encodeCursor(va[sz-1].CreatedAt)
	}

	pd := paginationData{
		pageCount:  len(va),
		totalCount: count,
		nextCursor: nextCursor,
		pager:      pager,
	}

	listResponse(c, va, pd)
}

func (r *Router) serverVersionedAttributesList(c *gin.Context) {
	pager, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "invalid pagination", err)
		return
	}

	srvUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		badRequestResponse(c, "invalid server uuid", err)
		return
	}

	dbVA, count, err := r.Store.ListVersionedAttributes(srvUUID, &pager)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	va := []VersionedAttributes{}

	for _, dbA := range dbVA {
		a := VersionedAttributes{}
		if err := a.fromDBModel(dbA); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		va = append(va, a)
	}

	nextCursor := ""

	sz := len(va)
	if sz != 0 {
		nextCursor = encodeCursor(va[sz-1].CreatedAt)
	}

	pd := paginationData{
		pageCount:  len(va),
		totalCount: count,
		nextCursor: nextCursor,
		pager:      pager,
	}

	listResponse(c, va, pd)
}

func (r *Router) serverVersionedAttributesCreate(c *gin.Context) {
	var va VersionedAttributes
	if err := c.ShouldBindJSON(&va); err != nil {
		badRequestResponse(c, "invalid versioned attributes", err)
		return
	}

	dbVA, err := va.toDBModel()
	if err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	srv, err := r.loadOrCreateServerFromParams(c)
	if err != nil {
		return
	}

	err = r.Store.CreateVersionedAttributes(srv, dbVA)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbVA.Namespace)
}
