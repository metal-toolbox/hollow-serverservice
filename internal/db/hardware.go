package db

import (
	"time"

	"github.com/google/uuid"
)

// Hardware represents a piece of hardware in a facility. These are the
// details of the physical hardware and are tracked separately from leases
// which track an instance of hardware.
type Hardware struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time `gorm:"index"`
	FacilityCode string
	BIOSConfigs  []BIOSConfig
}

// BeforeSave ensures that the hardware passes validation checks
// func (h *Hardware) BeforeSave(tx *gorm.DB) (err error) {
// 	if h.FacilityCode == "" {
// 		return requiredFieldMissing("hardware", "facility")
// 	}

// 	return nil
// }

// CreateHardware will persist hardware into the backend datastore
func CreateHardware(h Hardware) error {
	return db.Create(&h).Error
}

// HardwareList will return a list of hardware with the requested params
func HardwareList() ([]Hardware, error) {
	var hw []Hardware
	if err := db.Find(&hw).Error; err != nil {
		return nil, err
	}

	return hw, nil
}

// FindOrCreateHardwareByUUID will return an existing hardware instance if one
// already exists for the given UUID, if one doesn't then it will create a new
// instance in the database and return it.
func FindOrCreateHardwareByUUID(hwUUID uuid.UUID) (*Hardware, error) {
	var hw Hardware

	if err := db.FirstOrCreate(&hw, Hardware{ID: hwUUID}).Error; err != nil {
		return nil, err
	}

	return &hw, nil
}

// HardwareExists will return true or false based on if the UUID already exists
func HardwareExists(hwUUID uuid.UUID) bool {
	var count int64

	db.Model(&Hardware{}).Where("id = ?", hwUUID).Count(&count)

	return count == int64(1)
}
