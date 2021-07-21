package hollow

type pagination struct {
	Limit int    `json:"limit,omitempty"`
	Page  int    `json:"page,omitempty"`
	Sort  string `json:"sort,omitempty"`
}
