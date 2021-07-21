package db

import "gorm.io/gorm"

var (
	// MaxPaginationSize represents the maximum number of records that can be returned per page
	MaxPaginationSize = 1000
	// DefaultPaginationSize represents the default number of records that are returned per page
	DefaultPaginationSize = 100
)

// Pagination provides the parameters for paginating request
type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

func paginate(p Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := p.Page
		if page == 0 {
			page = 1
		}

		limit := p.Limit

		switch {
		case limit > MaxPaginationSize:
			limit = MaxPaginationSize
		case limit <= 0:
			limit = DefaultPaginationSize
		}

		offset := (page - 1) * limit

		if p.Sort != "" {
			db = db.Order(p.Sort)
		}

		return db.Offset(offset).Limit(limit)
	}
}
