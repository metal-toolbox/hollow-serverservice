package hollow

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.metalkube.net/hollow/internal/db"
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

	_, err = params.dbFilter(r)
	if err != nil {
		badRequestResponse(c, "invalid list params", err)
		return
	}

	dbSRV, err := db.Servers(
		qm.Load("Attributes"),
		qm.Load("VersionedAttributes"),
		qm.Load("ServerComponents.Attributes"),
		qm.Load("ServerComponents.ServerComponentType"),
	).All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count := int64(1)

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
	if err = srv.fromDBModel(dbSRV); err != nil {
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

	dbSRV, err := srv.toDBModel()
	if err != nil {
		badRequestResponse(c, "invalid server", err)
		return
	}

	if err := dbSRV.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbSRV.ID)
}

func (r *Router) serverDelete(c *gin.Context) {
	dbSRV, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	if _, err = dbSRV.Delete(c.Request.Context(), r.DB); err != nil {
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

	dbSRV, err := srv.toDBModel()
	if err != nil {
		badRequestResponse(c, "invalid server", err)
		return
	}

	cols := boil.Whitelist(db.ServerTableColumns.Name, db.ServerTableColumns.FacilityCode)

	if _, err := dbSRV.Update(c.Request.Context(), r.DB, cols); err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, u.String())
}

func (r *Router) serverVersionedAttributesGet(c *gin.Context) {
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	pager, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "invalid pagination", err)
		return
	}

	ns := c.Param("namespace")

	dbVA, err := srv.VersionedAttributes(db.VersionedAttributeWhere.Namespace.EQ(ns)).All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count, err := srv.VersionedAttributes(db.VersionedAttributeWhere.Namespace.EQ(ns)).Count(c.Request.Context(), r.DB)
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
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	pager, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "invalid pagination", err)
		return
	}

	dbVA, err := srv.VersionedAttributes().All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count := int64(1)

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
