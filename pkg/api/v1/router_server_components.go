package serverservice

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// serverComponentList returns a response with the list of components that matched the params.
func (r *Router) serverComponentList(c *gin.Context) {
	pager := parsePagination(c)

	params, err := parseQueryServerComponentsListParams(c)
	if err != nil {
		badRequestResponse(c, "invalid server component list params", err)
		return
	}

	dbSC, count, err := r.getServerComponents(c, params, pager)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	serverComponents := ServerComponentSlice{}

	for _, dbSC := range dbSC {
		sc := ServerComponent{}
		if err := sc.fromDBModel(dbSC); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		serverComponents = append(serverComponents, sc)
	}

	pd := paginationData{
		pageCount:  len(serverComponents),
		totalCount: count,
		pager:      pager,
	}

	listResponse(c, serverComponents, pd)
}

// serverComponentGet returns a response with the list of components referenced by the server UUID.
func (r *Router) serverComponentGet(c *gin.Context) {
	srv, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	pager := parsePagination(c)

	// list component by server_id
	// - include Attributes, VersionedAttributes and ServerComponentyType relations
	mods := []qm.QueryMod{
		qm.Load("Attributes"),
		qm.Load("VersionedAttributes",
			qm.Where(
				`(namespace, created_at) IN (SELECT namespace, MAX(created_at) FROM versioned_attributes WHERE server_id=? GROUP BY namespace, server_component_id)`,
				srv.ID,
			),
		),
		qm.Load("ServerComponentType"),
	}

	dbComps, err := srv.ServerComponents(mods...).All(c.Request.Context(), r.DB)
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

// serverComponentsCreate stores a ServerComponentSlice object into the backend store.
func (r *Router) serverComponentsCreate(c *gin.Context) {
	// load server based on the UUID parameter
	server, err := r.loadServerFromParams(c)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	// check server exists
	if server == nil {
		notFoundResponse(c, "no such server")
		return
	}

	// components payload
	var serverComponents ServerComponentSlice
	if err := c.ShouldBindJSON(&serverComponents); err != nil {
		badRequestResponse(c, "invalid payload: ServerComponentSlice", err)
		return
	}

	// component data is written in a transaction along with versioned attributes
	tx, err := r.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	// rollback is a no-op when the transaction is successful
	// nolint:errcheck // TODO(joel): log gerror instead of ignoring
	defer tx.Rollback()

	for _, srvComponent := range serverComponents {
		dbSrvComponent := srvComponent.toDBModel(server.ID)

		// insert component
		err := dbSrvComponent.Insert(c.Request.Context(), tx, boil.Infer())
		if err != nil {
			dbErrorResponse(c, err)
			return
		}

		// insert versioned attributes
		for _, versionedAttributes := range srvComponent.VersionedAttributes {
			dbVersionedAttributes := versionedAttributes.toDBModel()
			dbVersionedAttributes.ServerID = null.StringFrom(server.ID)
			dbVersionedAttributes.ServerComponentID = null.StringFrom(srvComponent.UUID.String())

			// insert true
			err = dbSrvComponent.AddVersionedAttributes(c.Request.Context(), tx, true, dbVersionedAttributes)
			if err != nil {
				dbErrorResponse(c, err)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, "")
}

// serverComponentUpdate updates existing server component attributes
func (r *Router) serverComponentUpdate(c *gin.Context) {
	// load server based on the UUID parameter
	server, err := r.loadServerFromParams(c)
	if err != nil {
		return
	}

	// check server exists
	if server == nil {
		notFoundResponse(c, "no such server")
		return
	}

	// components payload
	var serverComponents ServerComponentSlice
	if err := c.ShouldBindJSON(&serverComponents); err != nil {
		badRequestResponse(c, "invalid payload: ServerComponentSlice", err)
		return
	}

	// component data is written in a transaction along with versioned attributes
	tx, err := r.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	// rollback is a no-op when the transaction is successful
	// nolint:errcheck // TODO(joel): log gerror instead of ignoring
	defer tx.Rollback()

	for _, srvComponent := range serverComponents {
		dbSrvComponent := srvComponent.toDBModel(server.ID)

		// update component
		_, err := dbSrvComponent.Update(c.Request.Context(), r.DB, boil.Infer())
		if err != nil {
			dbErrorResponse(c, err)
			return
		}

		// update component versioned attributes
		for _, versionedAttributes := range srvComponent.VersionedAttributes {
			dbVersionedAttributes := versionedAttributes.toDBModel()
			dbVersionedAttributes.ServerID = null.StringFrom(server.ID)
			dbVersionedAttributes.ServerComponentID = null.StringFrom(srvComponent.UUID.String())

			err = dbSrvComponent.AddVersionedAttributes(c.Request.Context(), tx, true, dbVersionedAttributes)
			if err != nil {
				dbErrorResponse(c, err)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, "")
}
