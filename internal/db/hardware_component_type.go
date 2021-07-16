package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HardwareComponentType provides a way to group hardware components by the type
type HardwareComponentType struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"<-:create;not null;uniqueIndex;"`
}

// HardwareComponentTypeFilter provides the ability to filter to hardware that is returned for
// a query
type HardwareComponentTypeFilter struct {
	Name string
}

// BeforeSave ensures that the hardware component type passes validation checks
func (t *HardwareComponentType) BeforeSave(tx *gorm.DB) (err error) {
	if t.Name == "" {
		return requiredFieldMissing("hardware component type", "name")
	}

	return nil
}

// CreateHardwareComponentType will persist a hardware component type into the backend datastore
func (s *Store) CreateHardwareComponentType(t *HardwareComponentType) error {
	return s.db.Create(&t).Error
}

// GetHardwareComponentTypes will return a list of hardware component types with the requested params, if no
// filter is passed then it will return all hardware component types
func (s *Store) GetHardwareComponentTypes(filter *HardwareComponentTypeFilter) ([]HardwareComponentType, error) {
	var types []HardwareComponentType

	d := s.db

	if filter != nil {
		d = filter.apply(d)
	}

	if err := d.Find(&types).Error; err != nil {
		return nil, err
	}

	return types, nil
}

func (f *HardwareComponentTypeFilter) apply(d *gorm.DB) *gorm.DB {
	if f.Name != "" {
		d = d.Where("name = ?", f.Name)
	}

	return d
}
