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
	ID                  uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           time.Time `gorm:"index"`
	FacilityCode        string
	Attributes          []Attributes `gorm:"polymorphic:Entity;"`
	HardwareComponents  []HardwareComponent
	VersionedAttributes []VersionedAttributes `gorm:"polymorphic:Entity;"`
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
func (s *Store) CreateHardware(h *Hardware) error {
	return s.db.Create(h).Error
}

// GetHardware will return a list of hardware with the requested params, if no
// filter is passed then it will return all hardware
func (s *Store) GetHardware(filter *HardwareFilter) ([]Hardware, error) {
	var hw []Hardware

	d := hardwarePreload(s.db)

	if filter != nil {
		d = filter.apply(d)
	}

	if err := d.Find(&hw).Error; err != nil {
		return nil, err
	}

	return hw, nil
}

func hardwarePreload(db *gorm.DB) *gorm.DB {
	d := db.Preload("VersionedAttributes",
		"(created_at, namespace, entity_id, entity_type) IN (?)",
		db.Table("versioned_attributes").Select("max(created_at), namespace, entity_id, entity_type").Group("namespace").Group("entity_id").Group("entity_type"),
	)

	return d.Preload("HardwareComponents.HardwareComponentType").Preload("HardwareComponents.Attributes").Preload(clause.Associations)
}

// GetHardwareByUUID will return an existing hardware instance if one
// already exists for the given UUID.
func (s *Store) GetHardwareByUUID(hwUUID uuid.UUID) (*Hardware, error) {
	var hw Hardware

	err := hardwarePreload(s.db).First(&hw, Hardware{ID: hwUUID}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &hw, nil
}
