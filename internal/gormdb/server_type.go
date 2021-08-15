package gormdb

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// ServerType provides a way to group servers by their type
type ServerType struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Slug      string
}

// ServerTypeFilter provides the ability to filter the server components
// that are returned for a query
type ServerTypeFilter struct {
	Name string
	Slug string
}

// BeforeSave ensures that the server type passes validation checks
func (t *ServerType) BeforeSave(tx *gorm.DB) (err error) {
	if t.ID.String() == uuid.Nil.String() {
		t.ID = uuid.New()
	}

	if t.Name == "" {
		return requiredFieldMissing("server type", "name")
	}

	if t.Slug == "" {
		t.Slug = slug.Make(t.Name)
	}

	return nil
}

// CreateServerType will persist a server type into the backend datastore
func (s *Store) CreateServerType(t *ServerType) error {
	return s.db.Create(&t).Error
}

// GetServerTypes will return a list of server types with the requested params, if no
// filter is passed then it will return all server types
func (s *Store) GetServerTypes(filter *ServerTypeFilter, pager *Pagination) ([]ServerType, int64, error) {
	var (
		types []ServerType
		count int64
	)

	d := s.db

	if filter != nil {
		d = filter.apply(d)
	}

	if pager == nil {
		pager = &Pagination{}
	}

	if err := d.Scopes(paginate(*pager)).Find(&types).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return types, count, nil
}

func (f *ServerTypeFilter) apply(d *gorm.DB) *gorm.DB {
	if f.Name != "" {
		d = d.Where("name = ?", f.Name)
	}

	if f.Slug != "" {
		d = d.Where("name = ?", f.Slug)
	}

	return d
}
