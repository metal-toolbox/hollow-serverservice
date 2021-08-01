package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// ServerComponentType provides a way to group server components by their type
type ServerComponentType struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Slug      string
}

// ServerComponentTypeFilter provides the ability to filter the server components
// that are returned for a query
type ServerComponentTypeFilter struct {
	Name string
}

// BeforeSave ensures that the server component type passes validation checks
func (t *ServerComponentType) BeforeSave(tx *gorm.DB) (err error) {
	if t.ID.String() == uuid.Nil.String() {
		t.ID = uuid.New()
	}

	if t.Name == "" {
		return requiredFieldMissing("server component type", "name")
	}

	if t.Slug == "" {
		t.Slug = slug.Make(t.Name)
	}

	return nil
}

// CreateServerComponentType will persist a server component type into the backend datastore
func (s *Store) CreateServerComponentType(t *ServerComponentType) error {
	return s.db.Create(&t).Error
}

// GetServerComponentTypes will return a list of server component types with the requested params, if no
// filter is passed then it will return all server component types
func (s *Store) GetServerComponentTypes(filter *ServerComponentTypeFilter, pager *Pagination) ([]ServerComponentType, int64, error) {
	var (
		types []ServerComponentType
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

// FindServerComponentTypeBySlug will return a server component type with the matching slug
func (s *Store) FindServerComponentTypeBySlug(slug string) (*ServerComponentType, error) {
	var sct ServerComponentType

	err := s.db.First(&sct, ServerComponentType{Slug: slug}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &sct, nil
}

func (f *ServerComponentTypeFilter) apply(d *gorm.DB) *gorm.DB {
	if f.Name != "" {
		d = d.Where("name = ?", f.Name)
	}

	return d
}
