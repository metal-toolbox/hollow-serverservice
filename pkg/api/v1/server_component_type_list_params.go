package fleetdbapi

import "net/url"

// ServerComponentTypeListParams allows you to filter the results
type ServerComponentTypeListParams struct {
	Name             string
	PaginationParams *PaginationParams
}

func (p *ServerComponentTypeListParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	if p.Name != "" {
		q.Set("name", p.Name)
	}

	p.PaginationParams.setQuery(q)
}
