package db_test

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"go.metalkube.net/hollow/internal/db"
)

var (
	fixtureHardware      = db.Hardware{ID: uuid.New(), FacilityCode: "TEST1"}
	fixtureBIOSConfigOld = db.BIOSConfig{ID: uuid.New(), HardwareUUID: fixtureHardware.ID, ConfigValues: datatypes.JSON([]byte(`{"name": "old"}`))}
	fixtureBIOSConfigNew = db.BIOSConfig{ID: uuid.New(), HardwareUUID: fixtureHardware.ID, ConfigValues: datatypes.JSON([]byte(`{"name": "new"}`))}
)

func setupTestData() error {
	if err := db.CreateHardware(fixtureHardware); err != nil {
		return err
	}

	if err := db.CreateBIOSConfig(fixtureBIOSConfigOld); err != nil {
		return err
	}

	return db.CreateBIOSConfig(fixtureBIOSConfigNew)
}
