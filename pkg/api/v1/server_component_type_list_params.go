package hollow

import "net/url"

// ServerComponentTypeListParams allows you to filter the results
type ServerComponentTypeListParams struct {
	pagination
	Name string
}

func (f *ServerComponentTypeListParams) setQuery(q url.Values) {
	if f == nil {
		return
	}

	if f.Name != "" {
		q.Set("name", f.Name)
	}
}
