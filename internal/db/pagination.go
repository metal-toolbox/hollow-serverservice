package db

import (
	"time"

	"gorm.io/gorm"
)

var (
	// MaxPaginationSize represents the maximum number of records that can be returned per page
	MaxPaginationSize = 1000
	// DefaultPaginationSize represents the default number of records that are returned per page
	DefaultPaginationSize = 100
)

// Pagination provides the parameters for paginating request
type Pagination struct {
	Limit  int        `json:"limit"`
	Page   int        `json:"page"`
	Cursor *time.Time `json:"cursor"`
	// Sort  string `json:"sort"`
}

func paginate(p Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Order("created_at DESC").Limit(p.LimitUsed())

		switch {
		case p.Cursor != nil:
			db = db.Where("created_at < ?", p.Cursor)
		case p.Page != 0:
			db = db.Offset(p.Offset())
		}

		return db
	}
}

// LimitUsed returns the limit that was applied to the query
func (p *Pagination) LimitUsed() int {
	limit := p.Limit

	switch {
	case limit > MaxPaginationSize:
		limit = MaxPaginationSize
	case limit <= 0:
		limit = DefaultPaginationSize
	}

	return limit
}

// Offset returns the offset that was applied to the query
func (p *Pagination) Offset() int {
	page := p.Page
	if page == 0 {
		page = 1
	}

	return (page - 1) * p.LimitUsed()
}
