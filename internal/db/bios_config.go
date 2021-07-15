package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BIOSConfig represents the BIOS config of a given piece of hardware at a specific point in time
type BIOSConfig struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	HardwareID   uuid.UUID `gorm:"type:uuid;index"`
	Hardware     Hardware
	ConfigValues datatypes.JSON
	Timestamp    time.Time `json:"timestamp"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// BeforeSave ensures that the BIOS config passes validation checks
func (bc *BIOSConfig) BeforeSave(tx *gorm.DB) (err error) {
	if bc.HardwareID == uuid.Nil {
		return requiredFieldMissing("BIOSConfig", "hardware UUID")
	}

	return nil
}

// CreateBIOSConfig will persist a BIOSConfig into the backend datastore
func CreateBIOSConfig(bc *BIOSConfig) error {
	// if the hardware relation isn't already loaded ensure that the hardware does
	// exist. Otherwise return an error.
	if bc.HardwareID != uuid.Nil {
		if !HardwareExists(bc.HardwareID) {
			return newNotFoundError("hardware UUID")
		}
	}

	return db.Create(bc).Error
}

// BIOSConfigList will return all the BIOSConfigs for a given Hardware UUID, the list will be sorted with the newest one
// first
func BIOSConfigList(hwUUID uuid.UUID) ([]BIOSConfig, error) {
	var bcl []BIOSConfig
	if err := db.Where(&BIOSConfig{HardwareID: hwUUID}).Order("created_at desc").Find(&bcl).Error; err != nil {
		return nil, err
	}

	return bcl, nil
}
