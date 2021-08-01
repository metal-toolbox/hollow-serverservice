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

	FixtureSCTFins = ServerComponentType{ID: uuid.New(), Name: "Fins", Slug: "fins"}

	FixtureAttributesNemoMetadata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceMetadata, Data: datatypes.JSON([]byte(`{"age":6,"location":"Fishbowl"}`))}
	FixtureAttributesDoryMetadata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceMetadata, Data: datatypes.JSON([]byte(`{"age":12,"location":"East Austalian Current"}`))}
	FixtureAttributesMarlinMetadata = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceMetadata, Data: datatypes.JSON([]byte(`{"age":10,"location":"East Austalian Current"}`))}

	FixtureAttributesNemoOtherdata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceOtherdata, Data: datatypes.JSON([]byte(`{"enabled": true, "type": "clown", "lastUpdated": 1624960800, "nested": {"tag": "finding-nemo", "number": 1}}`))}
	FixtureAttributesDoryOtherdata   = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceOtherdata, Data: datatypes.JSON([]byte(`{"enabled": true, "type": "blue-tang", "lastUpdated": 1624960400, "nested": {"tag": "finding-nemo", "number": 2}}`))}
	FixtureAttributesMarlinOtherdata = Attributes{ID: uuid.New(), Namespace: FixtureNamespaceOtherdata, Data: datatypes.JSON([]byte(`{"enabled": false, "type": "clown", "lastUpdated": 1624960000, "nested": {"tag": "finding-nemo", "number": 3}}`))}

	FixtureSCDoryLeftFin    = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "Normal Fin", Serial: "Left"}
	FixtureSCDoryRightFin   = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "Normal Fin", Serial: "Right"}
	FixtureSCMarlinLeftFin  = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "Normal Fin", Serial: "Left"}
	FixtureSCMarlinRightFin = ServerComponent{ID: uuid.New(), ServerComponentTypeID: FixtureSCTFins.ID, Model: "Normal Fin", Serial: "Right"}

	FixtureVersionedAttributesNemoRightFin = VersionedAttributes{
		ID:        uuid.New(),
		Namespace: FixtureNamespaceVersioned,
		Data:      datatypes.JSON([]byte(`{"something": "cool"}`)),
	}

	FixtureSCNemoLeftFin = ServerComponent{
		ID:                    uuid.New(),
		ServerComponentTypeID: FixtureSCTFins.ID,
		Model:                 "Normal Fin",
		Serial:                "Left",
	}

	FixtureSCNemoRightFin = ServerComponent{
		ID:                    uuid.New(),
		ServerComponentTypeID: FixtureSCTFins.ID,
		Name:                  "My Lucky Fin",
		Vendor:                "Barracuda",
		Model:                 "A Lucky Fin",
		Serial:                "Right",
		VersionedAttributes: []VersionedAttributes{
			FixtureVersionedAttributesNemoRightFin,
		},
	}

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
		ID:        uuid.New(),
		ServerID:  &FixtureServerNemo.ID,
		Namespace: FixtureNamespaceVersioned,
		Data:      datatypes.JSON([]byte(`{"name": "old"}`)),
	}

	FixtureVersionedAttributesNew = VersionedAttributes{
		ID:        uuid.New(),
		ServerID:  &FixtureServerNemo.ID,
		Namespace: FixtureNamespaceVersioned,
		Data:      datatypes.JSON([]byte(`{"name": "new"}`)),
	}

	FixtureServer = []Server{FixtureServerNemo, FixtureServerDory, FixtureServerMarlin}
)

func (s *Store) setupTestData() error {
	if err := s.CreateServerComponentType(&FixtureSCTFins); err != nil {
		return err
	}

	for _, srv := range FixtureServer {
		if err := s.CreateServer(&srv); err != nil {
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
