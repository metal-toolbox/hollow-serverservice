package db_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateAttributes(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		a           *db.Attributes
		expectError bool
		errorMsg    string
	}{
		{"missing namespace", &db.Attributes{}, true, "validation failed: namespace is a required attributes attribute"},
		{"happy path", &db.Attributes{ServerID: &db.FixtureServerDory.ID, Namespace: "hollow.test", Data: datatypes.JSON([]byte(`{"value": "set"}`))}, false, ""},
	}

	for _, tt := range testCases {
		err := s.CreateAttributes(tt.a)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
		}
	}
}

func TestDeleteAttributes(t *testing.T) {
	s := db.DatabaseTest(t)

	err := s.DeleteAttributes(&db.FixtureAttributesDoryMetadata)
	assert.NoError(t, err)
}

func TestGetAttributesByServerUUID(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName         string
		u                uuid.UUID
		expectUUIDs      []uuid.UUID
		expectedNotFound bool
	}{
		{
			"happy path",
			db.FixtureServerDory.ID,
			[]uuid.UUID{
				db.FixtureAttributesDoryMetadata.ID,
				db.FixtureAttributesDoryOtherdata.ID,
			},
			false,
		},
		{
			"not found server uuid",
			uuid.New(),
			[]uuid.UUID{},
			true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			attrs, count, err := s.GetAttributesByServerUUID(tt.u, nil)
			if tt.expectedNotFound {
				assert.Error(t, err)
				assert.ErrorIs(t, err, db.ErrNotFound)
			} else {
				assert.NoError(t, err)
				assert.Len(t, attrs, int(count))

				var respIDs []uuid.UUID
				for _, a := range attrs {
					respIDs = append(respIDs, a.ID)
				}

				assert.ElementsMatch(t, tt.expectUUIDs, respIDs)
			}
		})
	}
}

func TestGetAttributesByServerUUIDAndNamespace(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName         string
		u                uuid.UUID
		ns               string
		expectedID       uuid.UUID
		expectedNotFound bool
	}{
		{"happy path", db.FixtureServerDory.ID, db.FixtureNamespaceMetadata, db.FixtureAttributesDoryMetadata.ID, false},
		{"not found server uuid", uuid.New(), db.FixtureNamespaceMetadata, uuid.Nil, true},
		{"not found namespace", db.FixtureServerDory.ID, "unknown", uuid.Nil, true},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			attr, err := s.GetAttributesByServerUUIDAndNamespace(tt.u, tt.ns)

			if tt.expectedNotFound {
				assert.Error(t, err)
				assert.ErrorIs(t, err, db.ErrNotFound)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, attr)
				assert.Equal(t, tt.expectedID, attr.ID)
			}
		})
	}
}

func TestUpdateAttributesByServerUUIDAndNamespace(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName         string
		u                uuid.UUID
		ns               string
		data             json.RawMessage
		expectedNotFound bool
	}{
		{"happy path", db.FixtureServerDory.ID, db.FixtureNamespaceMetadata, json.RawMessage([]byte(`{"age": 12, "location": "Fishbowl"}`)), false},
		{"not found server uuid", uuid.New(), db.FixtureNamespaceMetadata, json.RawMessage([]byte(`{"age": 12, "location": "Fishbowl"}`)), true},
		{"happy path - new namespace should upsert", db.FixtureServerDory.ID, "unknown", json.RawMessage([]byte(`{"age": 12, "location": "Fishbowl"}`)), false},
		{"no namespace provided", db.FixtureServerDory.ID, "", json.RawMessage{}, true},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			err := s.UpdateAttributesByServerUUIDAndNamespace(tt.u, tt.ns, tt.data)

			if tt.expectedNotFound {
				assert.Error(t, err)
				assert.ErrorIs(t, err, db.ErrNotFound)
			} else {
				assert.NoError(t, err)
				attr, err := s.GetAttributesByServerUUIDAndNamespace(tt.u, tt.ns)
				assert.NoError(t, err)
				assert.Equal(t, datatypes.JSON(tt.data), attr.Data)
			}
		})
	}
}
