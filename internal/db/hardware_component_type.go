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
	Name      string `gorm:"<-:create;not null;"`
}

// BeforeSave ensures that the hardware component type passes validation checks
func (t *HardwareComponentType) BeforeSave(tx *gorm.DB) (err error) {
	if t.Name == "" {
		return requiredFieldMissing("hardware component type", "name")
	}

	return nil
}

// CreateHardwareComponentType will persist a hardware component type into the backend datastore
func CreateHardwareComponentType(t HardwareComponentType) error {
	return db.Create(&t).Error
}
