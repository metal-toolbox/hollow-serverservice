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
	FixtureNamespaceVersioned = "hollow.versioned"

	FixtureSCTFins = ServerComponentType{ID: uuid.New(), Name: "Fins"}

	FixtureAttributesNemoMetadata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceMetadata, Values: datatypes.JSON([]byte(`{"location": "Fishbowl", "age": 6}`))}
	FixtureAttributesDoryMetadata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceMetadata, Values: datatypes.JSON([]byte(`{"location": "East Austalian Current", "age": 12}`))}
	FixtureAttributesMarlinMetadata = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceMetadata, Values: datatypes.JSON([]byte(`{"location": "East Austalian Current", "age": 10}`))}

	FixtureAttributesNemoOtherdata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceOtherdata, Values: datatypes.JSON([]byte(`{"enabled": true, "type": "clown", "lastUpdated": 1624960800, "nested": {"tag": "finding-nemo", "number": 1}}`))}
	FixtureAttributesDoryOtherdata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceOtherdata, Values: datatypes.JSON([]byte(`{"enabled": true, "type": "blue-tang", "lastUpdated": 1624960400, "nested": {"tag": "finding-nemo", "number": 2}}`))}
	FixtureAttributesMarlinOtherdata = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceOtherdata, Values: datatypes.JSON([]byte(`{"enabled": false, "type": "clown", "lastUpdated": 1624960000, "nested": {"tag": "finding-nemo", "number": 3}}`))}

	FixtureSCNemoLeftFin    = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "Normal Fin", Serial: "Left"}
	FixtureSCNemoRightFin   = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "A Lucky Fin", Serial: "Right"}
	FixtureSCDoryLeftFin    = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "Normal Fin", Serial: "Left"}
	FixtureSCDoryRightFin   = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "Normal Fin", Serial: "Right"}
	FixtureSCMarlinLeftFin  = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "Normal Fin", Serial: "Left"}
	FixtureSCMarlinRightFin = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "Normal Fin", Serial: "Right"}

	FixtureServerNemo = Server{
		ID:               uuid.New(),
		Name:             "Nemo",
		FacilityCode:     "Sydney",
		Attributes:       []Attributes{FixtureAttributesNemoMetadata, FixtureAttributesNemoOtherdata},
		ServerComponents: []ServerComponent{FixtureSCNemoLeftFin, FixtureSCNemoRightFin},
	}

	FixtureServerDory = Server{
		ID:               uuid.New(),
		Name:             "Dory",
		FacilityCode:     "Ocean",
		Attributes:       []Attributes{FixtureAttributesDoryMetadata, FixtureAttributesDoryOtherdata},
		ServerComponents: []ServerComponent{FixtureSCDoryLeftFin, FixtureSCDoryRightFin},
	}

	FixtureServerMarlin = Server{
		ID:               uuid.New(),
		Name:             "Marlin",
		FacilityCode:     "Ocean",
		Attributes:       []Attributes{FixtureAttributesMarlinMetadata, FixtureAttributesMarlinOtherdata},
		ServerComponents: []ServerComponent{FixtureSCMarlinLeftFin, FixtureSCMarlinRightFin},
	}

	FixtureVersionedAttributesOld = VersionedAttributes{
		ID:         uuid.New(),
		EntityType: "servers",
		EntityID:   FixtureServerNemo.ID,
		Namespace:  FixtureNamespaceVersioned,
		Values:     datatypes.JSON([]byte(`{"name": "old"}`)),
	}

	FixtureVersionedAttributesNew = VersionedAttributes{
		ID:         uuid.New(),
		EntityType: "servers",
		EntityID:   FixtureServerNemo.ID,
		Namespace:  FixtureNamespaceVersioned,
		Values:     datatypes.JSON([]byte(`{"name": "new"}`)),
	}

	FixtureServer = []Server{FixtureServerNemo, FixtureServerDory, FixtureServerMarlin}
)

func (s *Store) setupTestData() error {
	if err := s.CreateServerComponentType(&FixtureSCTFins); err != nil {
		return err
	}

	for _, hw := range FixtureServer {
		if err := s.CreateServer(&hw); err != nil {
			return err
		}
	}

	for _, a := range []VersionedAttributes{FixtureVersionedAttributesOld, FixtureVersionedAttributesNew} {
		if err := s.db.Create(&a).Error; err != nil {
			return err
		}
	}

	return nil
}
