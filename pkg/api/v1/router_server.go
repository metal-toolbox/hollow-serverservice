package hollow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Router) serverList(c *gin.Context) {
	pager := parsePagination(c)

	var params ServerListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid filter",
			"error":   err.Error(),
		})
	}

	alp, err := parseQueryAttributesListParams(c, "attr")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid attributes list params",
			"error":   err.Error(),
		})
	}

	params.AttributeListParams = alp

	valp, err := parseQueryAttributesListParams(c, "ver_attr")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid versioned attributes list params",
			"error":   err.Error(),
		})
	}

	params.VersionedAttributeListParams = valp

	sclp, err := parseQueryServerComponentsListParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid server component list params",
			"error":   err.Error(),
		})
	}

	params.ComponentListParams = sclp

	dbFilter, err := params.dbFilter()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid list params",
			"error":   err.Error(),
		})
	}

	dbSRV, err := r.Store.GetServers(dbFilter, &pager)
	if err != nil {
		dbFailureResponse(c, err)
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

	c.JSON(http.StatusOK, srvs)
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

	c.JSON(http.StatusOK, srv)
}

func (r *Router) serverCreate(c *gin.Context) {
	var srv Server
	if err := c.ShouldBindJSON(&srv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid server",
			"error":   err.Error(),
		})

		return
	}

	dbSRV, err := srv.toDBModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid server", "error": err.Error})
		return
	}

	if err := r.Store.CreateServer(dbSRV); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed to create server", err))
		return
	}

	createdResponse(c, &dbSRV.ID)
}

func (r *Router) serverDelete(c *gin.Context) {
	dbSRV, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	if err = r.Store.DeleteServer(dbSRV); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed deleting resource", err))
		return
	}

	deletedResponse(c)
}

func (r *Router) serverVersionedAttributesList(c *gin.Context) {
	srvUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid server uuid", "error": err.Error()})
		return
	}

	dbVA, err := r.Store.GetVersionedAttributes(srvUUID)
	if err != nil {
		dbFailureResponse(c, err)
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

	c.JSON(http.StatusOK, va)
}

func (r *Router) serverVersionedAttributesCreate(c *gin.Context) {
	var va VersionedAttributes
	if err := c.ShouldBindJSON(&va); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid versioned attributes",
			"error":   err.Error(),
		})

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
		dbFailureResponse(c, err)
		return
	}

	createdResponse(c, &dbVA.ID)
}
