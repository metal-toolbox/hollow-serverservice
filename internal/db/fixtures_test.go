package db_test

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"go.metalkube.net/hollow/internal/db"
)

var (
	fixtureNamespaceMetadata  = "hollow.metadata"
	fixtureNamespaceOtherdata = "hollow.other_data"

	fixtureAttributesNemoMetadata   = db.Attributes{ID: uuid.New(), Namespace: fixtureNamespaceMetadata, Values: datatypes.JSON([]byte(`{"location": "Fishbowl", "age": 6}`))}
	fixtureAttributesDoryMetadata   = db.Attributes{ID: uuid.New(), Namespace: fixtureNamespaceMetadata, Values: datatypes.JSON([]byte(`{"location": "East Austalian Current", "age": 10}`))}
	fixtureAttributesMarlinMetadata = db.Attributes{ID: uuid.New(), Namespace: fixtureNamespaceMetadata, Values: datatypes.JSON([]byte(`{"location": "East Austalian Current", "age": 10}`))}

	fixtureAttributesNemoOtherdata   = db.Attributes{ID: uuid.New(), Namespace: fixtureNamespaceOtherdata, Values: datatypes.JSON([]byte(`{"enabled": true, "type": "clown", "lastUpdated": 1624960800, "nested": {"tag": "finding-nemo", "number": 1}}`))}
	fixtureAttributesDoryOtherdata   = db.Attributes{ID: uuid.New(), Namespace: fixtureNamespaceOtherdata, Values: datatypes.JSON([]byte(`{"enabled": true, "type": "blue-tang", "lastUpdated": 1624960400, "nested": {"tag": "finding-nemo", "number": 2}}`))}
	fixtureAttributesMarlinOtherdata = db.Attributes{ID: uuid.New(), Namespace: fixtureNamespaceOtherdata, Values: datatypes.JSON([]byte(`{"enabled": false, "type": "clown", "lastUpdated": 1624960000, "nested": {"tag": "finding-nemo", "number": 3}}`))}

	fixtureHardwareNemo   = db.Hardware{ID: uuid.New(), FacilityCode: "Nemo", Attributes: []db.Attributes{fixtureAttributesNemoMetadata, fixtureAttributesNemoOtherdata}}
	fixtureHardwareDory   = db.Hardware{ID: uuid.New(), FacilityCode: "Dory", Attributes: []db.Attributes{fixtureAttributesDoryMetadata, fixtureAttributesDoryOtherdata}}
	fixtureHardwareMarlin = db.Hardware{ID: uuid.New(), FacilityCode: "Marlin", Attributes: []db.Attributes{fixtureAttributesMarlinMetadata, fixtureAttributesMarlinOtherdata}}

	fixtureBIOSConfig    = db.BIOSConfig{ID: uuid.New(), HardwareID: fixtureHardwareNemo.ID, ConfigValues: datatypes.JSON([]byte(`{"name": "old"}`))}
	fixtureBIOSConfigNew = db.BIOSConfig{ID: uuid.New(), HardwareID: fixtureHardwareNemo.ID, ConfigValues: datatypes.JSON([]byte(`{"name": "new"}`))}
)

func setupTestData() error {
	for _, hw := range []db.Hardware{fixtureHardwareNemo, fixtureHardwareDory, fixtureHardwareMarlin} {
		if err := db.CreateHardware(hw); err != nil {
			return err
		}
	}

	if err := db.CreateBIOSConfig(fixtureBIOSConfig); err != nil {
		return err
	}

	return db.CreateBIOSConfig(fixtureBIOSConfigNew)
}
