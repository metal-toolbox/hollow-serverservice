package hollow

import "net/url"

// HardwareComponentTypeListParams allows you to filter the results
type HardwareComponentTypeListParams struct {
	pagination
	Name string
}

func (f *HardwareComponentTypeListParams) setQuery(q url.Values) {
	if f == nil {
		return
	}

	if f.Name != "" {
		q.Set("name", f.Name)
	}
}
