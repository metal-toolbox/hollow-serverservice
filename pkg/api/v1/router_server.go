package dcim

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"

	"go.hollow.sh/serverservice/internal/models"
)

func (r *Router) serverList(c *gin.Context) {
	pager := parsePagination(c)

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

	params.PaginationParams = &pager

	dbSRV, count, err := r.getServers(c, params)
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

	pd := paginationData{
		pageCount:  len(srvs),
		totalCount: count,
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
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	var newValues Server
	if err := c.ShouldBindJSON(&newValues); err != nil {
		badRequestResponse(c, "invalid server", err)
		return
	}

	srv.Name = null.StringFrom(newValues.Name)
	srv.FacilityCode = null.StringFrom(newValues.FacilityCode)

	cols := boil.Infer()

	if _, err := srv.Update(c.Request.Context(), r.DB, cols); err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, srv.ID)
}

func (r *Router) serverVersionedAttributesGet(c *gin.Context) {
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	pager := parsePagination(c)

	ns := c.Param("namespace")

	dbVA, err := srv.VersionedAttributes(models.VersionedAttributeWhere.Namespace.EQ(ns), qm.OrderBy("created_at DESC")).All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count, err := srv.VersionedAttributes(models.VersionedAttributeWhere.Namespace.EQ(ns)).Count(c.Request.Context(), r.DB)
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

	pd := paginationData{
		pageCount:  len(va),
		totalCount: count,
		pager:      pager,
	}

	listResponse(c, va, pd)
}

func (r *Router) serverVersionedAttributesList(c *gin.Context) {
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	pager := parsePagination(c)

	dbVA, err := srv.VersionedAttributes(qm.OrderBy("created_at DESC")).All(c.Request.Context(), r.DB)
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

	count := int64(len(va))

	pd := paginationData{
		pageCount:  len(va),
		totalCount: count,
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

	dbVA := va.toDBModel()

	srv, err := r.loadOrCreateServerFromParams(c)
	if err != nil {
		return
	}

	// nolint:errcheck If this fails continue on
	curVA, _ := srv.VersionedAttributes(qm.Where("namespace = ?", va.Namespace), qm.OrderBy("created_at DESC")).One(c.Request.Context(), r.DB)

	if curVA != nil && areEqualJSON(dbVA.Data, curVA.Data) {
		curVA.Tally++

		_, err := curVA.Update(c.Request.Context(), r.DB, boil.Whitelist("tally", "updated_at"))
		if err != nil {
			dbErrorResponse(c, err)
			return
		}

		createdResponse(c, curVA.Namespace)

		return
	}

	if err := srv.AddVersionedAttributes(c.Request.Context(), r.DB, true, dbVA); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbVA.Namespace)
}

func areEqualJSON(s1, s2 types.JSON) bool {
	var (
		o1 interface{}
		o2 interface{}
	)

	if err := json.Unmarshal([]byte(s1), &o1); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(s2), &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
