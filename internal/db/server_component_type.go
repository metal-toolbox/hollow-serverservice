package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServerComponentType provides a way to group server components by their type
type ServerComponentType struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"<-:create;not null;uniqueIndex;"`
}

// ServerComponentTypeFilter provides the ability to filter the server components
// that are returned for a query
type ServerComponentTypeFilter struct {
	Name string
}

// BeforeSave ensures that the server component type passes validation checks
func (t *ServerComponentType) BeforeSave(tx *gorm.DB) (err error) {
	if t.Name == "" {
		return requiredFieldMissing("server component type", "name")
	}

	return nil
}

// CreateServerComponentType will persist a server component type into the backend datastore
func (s *Store) CreateServerComponentType(t *ServerComponentType) error {
	return s.db.Create(&t).Error
}

// GetServerComponentTypes will return a list of server component types with the requested params, if no
// filter is passed then it will return all server component types
func (s *Store) GetServerComponentTypes(filter *ServerComponentTypeFilter, pager *Pagination) ([]ServerComponentType, error) {
	var types []ServerComponentType

	d := s.db

	if filter != nil {
		d = filter.apply(d)
	}

	if pager == nil {
		pager = &Pagination{}
	}

	if err := d.Scopes(paginate(*pager)).Find(&types).Error; err != nil {
		return nil, err
	}

	return types, nil
}

func (f *ServerComponentTypeFilter) apply(d *gorm.DB) *gorm.DB {
	if f.Name != "" {
		d = d.Where("name = ?", f.Name)
	}

	return d
}
