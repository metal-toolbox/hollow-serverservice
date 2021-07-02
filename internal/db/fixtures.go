//+build testtools

package db

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

//nolint:revive
var (
	FixtureNamespaceMetadata  = "hollow.metadata"
	FixtureNamespaceOtherdata = "hollow.other_data"

	FixtureHCTFins = HardwareComponentType{ID: uuid.New(), Name: "Fins"}

	FixtureAttributesNemoMetadata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceMetadata, Values: datatypes.JSON([]byte(`{"location": "Fishbowl", "age": 6}`))}
	FixtureAttributesDoryMetadata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceMetadata, Values: datatypes.JSON([]byte(`{"location": "East Austalian Current", "age": 10}`))}
	FixtureAttributesMarlinMetadata = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceMetadata, Values: datatypes.JSON([]byte(`{"location": "East Austalian Current", "age": 10}`))}

	FixtureAttributesNemoOtherdata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceOtherdata, Values: datatypes.JSON([]byte(`{"enabled": true, "type": "clown", "lastUpdated": 1624960800, "nested": {"tag": "finding-nemo", "number": 1}}`))}
	FixtureAttributesDoryOtherdata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceOtherdata, Values: datatypes.JSON([]byte(`{"enabled": true, "type": "blue-tang", "lastUpdated": 1624960400, "nested": {"tag": "finding-nemo", "number": 2}}`))}
	FixtureAttributesMarlinOtherdata = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceOtherdata, Values: datatypes.JSON([]byte(`{"enabled": false, "type": "clown", "lastUpdated": 1624960000, "nested": {"tag": "finding-nemo", "number": 3}}`))}

	FixtureHCNemoLeftFin    = HardwareComponent{ID: uuid.New(), HardwareComponentTypeID: FixtureHCTFins.ID, Model: "Normal Fin", Serial: "Left"}
	FixtureHCNemoRightFin   = HardwareComponent{ID: uuid.New(), HardwareComponentTypeID: FixtureHCTFins.ID, Model: "A Lucky Fin", Serial: "Right"}
	FixtureHCDoryLeftFin    = HardwareComponent{ID: uuid.New(), HardwareComponentTypeID: FixtureHCTFins.ID, Model: "Normal Fin", Serial: "Left"}
	FixtureHCDoryRightFin   = HardwareComponent{ID: uuid.New(), HardwareComponentTypeID: FixtureHCTFins.ID, Model: "Normal Fin", Serial: "Right"}
	FixtureHCMarlinLeftFin  = HardwareComponent{ID: uuid.New(), HardwareComponentTypeID: FixtureHCTFins.ID, Model: "Normal Fin", Serial: "Left"}
	FixtureHCMarlinRightFin = HardwareComponent{ID: uuid.New(), HardwareComponentTypeID: FixtureHCTFins.ID, Model: "Normal Fin", Serial: "Right"}

	FixtureHardwareNemo = Hardware{
		ID:                 uuid.New(),
		FacilityCode:       "Nemo",
		Attributes:         []Attributes{FixtureAttributesNemoMetadata, FixtureAttributesNemoOtherdata},
		HardwareComponents: []HardwareComponent{FixtureHCNemoLeftFin, FixtureHCNemoRightFin},
	}

	FixtureHardwareDory = Hardware{
		ID:                 uuid.New(),
		FacilityCode:       "Dory",
		Attributes:         []Attributes{FixtureAttributesDoryMetadata, FixtureAttributesDoryOtherdata},
		HardwareComponents: []HardwareComponent{FixtureHCDoryLeftFin, FixtureHCDoryRightFin},
	}

	FixtureHardwareMarlin = Hardware{
		ID:                 uuid.New(),
		FacilityCode:       "Marlin",
		Attributes:         []Attributes{FixtureAttributesMarlinMetadata, FixtureAttributesMarlinOtherdata},
		HardwareComponents: []HardwareComponent{FixtureHCMarlinLeftFin, FixtureHCMarlinRightFin},
	}

	FixtureBIOSConfig    = BIOSConfig{ID: uuid.New(), HardwareID: FixtureHardwareNemo.ID, ConfigValues: datatypes.JSON([]byte(`{"name": "old"}`))}
	FixtureBIOSConfigNew = BIOSConfig{ID: uuid.New(), HardwareID: FixtureHardwareNemo.ID, ConfigValues: datatypes.JSON([]byte(`{"name": "new"}`))}
)

func setupTestData() error {
	if err := CreateHardwareComponentType(FixtureHCTFins); err != nil {
		return err
	}

	for _, hw := range []Hardware{FixtureHardwareNemo, FixtureHardwareDory, FixtureHardwareMarlin} {
		if err := CreateHardware(hw); err != nil {
			return err
		}
	}

	for _, bc := range []BIOSConfig{FixtureBIOSConfig, FixtureBIOSConfigNew} {
		if err := CreateBIOSConfig(bc); err != nil {
			return err
		}
	}

	return nil
}
