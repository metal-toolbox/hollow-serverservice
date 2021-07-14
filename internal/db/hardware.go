package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Hardware represents a piece of hardware in a facility. These are the
// details of the physical hardware and are tracked separately from leases
// which track an instance of hardware.
type Hardware struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          time.Time `gorm:"index"`
	FacilityCode       string
	BIOSConfigs        []BIOSConfig
	Attributes         []Attributes `gorm:"polymorphic:Entity;"`
	HardwareComponents []HardwareComponent
}

// HardwareFilter provides the ability to filter to hardware that is returned for
// a query
type HardwareFilter struct {
	FacilityCode      string
	AttributesFilters []AttributesFilter
}

// TableName overrides the table name used by Hardware to `hardware`
func (Hardware) TableName() string {
	return "hardware"
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

// GetHardware will return a list of hardware with the requested params, if no
// filter is passed then it will return all hardware
func GetHardware(filter *HardwareFilter) ([]Hardware, error) {
	var hw []Hardware

	d := db

	if filter != nil {
		d = filter.apply(db)
	}

	if err := hardwarePreload(d).Find(&hw).Error; err != nil {
		return nil, err
	}

	return hw, nil
}

func hardwarePreload(db *gorm.DB) *gorm.DB {
	return db.Preload("HardwareComponents.HardwareComponentType").Preload("HardwareComponents.Attributes").Preload(clause.Associations)
}

// FindOrCreateHardwareByUUID will return an existing hardware instance if one
// already exists for the given UUID, if one doesn't then it will create a new
// instance in the database and return it.
func FindOrCreateHardwareByUUID(hwUUID uuid.UUID) (*Hardware, error) {
	var hw Hardware

	if err := hardwarePreload(db).FirstOrCreate(&hw, Hardware{ID: hwUUID}).Error; err != nil {
		return nil, err
	}

	return &hw, nil
}

// FindHardwareByUUID will return an existing hardware instance if one
// already exists for the given UUID.
func FindHardwareByUUID(hwUUID uuid.UUID) (*Hardware, error) {
	var hw Hardware

	err := hardwarePreload(db).First(&hw, Hardware{ID: hwUUID}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
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

func (f *HardwareFilter) apply(d *gorm.DB) *gorm.DB {
	if f.FacilityCode != "" {
		d = d.Where("facility_code = ?", f.FacilityCode)
	}

	if f.AttributesFilters != nil {
		for i, af := range f.AttributesFilters {
			d = af.apply(d, i)
		}
	}

	return d
}
