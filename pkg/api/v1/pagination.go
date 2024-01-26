package serverservice

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

var (
	// maxPaginationSize represents the maximum number of records that can be returned per page
	maxPaginationSize = 1000
	// defaultPaginationSize represents the default number of records that are returned per page
	defaultPaginationSize = 100
)

// PaginationParams allow you to paginate the results
type PaginationParams struct {
	Limit   int    `json:"limit,omitempty"`
	Page    int    `json:"page,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
	Preload bool   `json:"preload,omitempty"`
	OrderBy string `json:"orderby,omitempty"`
}

type paginationData struct {
	pageCount  int
	totalCount int64
	pager      PaginationParams
}

func parsePagination(c *gin.Context) PaginationParams {
	// Initializing default
	limit := defaultPaginationSize
	page := 1
	query := c.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		}
	}

	return PaginationParams{
		Limit: limit,
		Page:  page,
	}
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (p *PaginationParams) queryMods() []qm.QueryMod {
	if p == nil {
		p = &PaginationParams{}
	}

	mods := []qm.QueryMod{}

	mods = append(mods, qm.Limit(p.limitUsed()))

	if p.Page != 0 {
		mods = append(mods, qm.Offset(p.offset()))
	}

	// match the old functionality for now...will handle order and load as params later
	if p.OrderBy != "" {
		mods = append(mods, qm.OrderBy(p.OrderBy))
	}

	return mods
}

// serverQueryMods queryMods converts the list params into sql conditions that can be added to sql queries
func (p *PaginationParams) serverQueryMods() []qm.QueryMod {
	if p == nil {
		p = &PaginationParams{}
	}

	mods := []qm.QueryMod{}

	mods = append(mods, qm.Limit(p.limitUsed()))

	if p.Page != 0 {
		mods = append(mods, qm.Offset(p.offset()))
	}

	// match the old functionality for now...will handle order and load as params later
	if p.OrderBy != "" {
		mods = append(mods, qm.OrderBy(p.OrderBy))
	}

	if p.Preload {
		preload := []qm.QueryMod{
			qm.Load("Attributes"),
			qm.Load("VersionedAttributes", qm.Where("(server_id, namespace, created_at) IN (select server_id, namespace, max(created_at) from versioned_attributes group by server_id, namespace)")),
			qm.Load("ServerComponents.Attributes"),
			qm.Load("ServerComponents.ServerComponentType"),
		}
		mods = append(mods, preload...)
	}

	return mods
}

// serverComponentQueryMods converts the server component list params into sql conditions that can be added to sql queries
func (p *PaginationParams) serverComponentsQueryMods() []qm.QueryMod {
	if p == nil {
		p = &PaginationParams{}
	}

	mods := []qm.QueryMod{}

	mods = append(mods, qm.Limit(p.limitUsed()))

	if p.Page != 0 {
		mods = append(mods, qm.Offset(p.offset()))
	}

	mods = append(mods, qm.OrderBy(models.ServerComponentTableColumns.CreatedAt+" DESC"))

	preload := []qm.QueryMod{
		qm.Load("Attributes"),
		qm.Load("VersionedAttributes", qm.Where("(server_component_id, namespace, created_at) IN (select server_component_id, namespace, max(created_at) from versioned_attributes group by server_component_id, namespace)")),
		qm.Load("ServerComponentType"),
	}
	mods = append(mods, preload...)

	return mods
}

func (p *PaginationParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}

	if p.Page != 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}

	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
}

func (p *PaginationParams) limitUsed() int {
	limit := p.Limit

	switch {
	case limit > maxPaginationSize:
		limit = maxPaginationSize
	case limit <= 0:
		limit = defaultPaginationSize
	}

	return limit
}

func (p *PaginationParams) offset() int {
	page := p.Page
	if page == 0 {
		page = 1
	}

	return (page - 1) * p.limitUsed()
}
