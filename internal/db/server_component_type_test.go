package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateServerComponentType(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		ct          *db.ServerComponentType
		expectError bool
		errorMsg    string
	}{
		{"missing name", &db.ServerComponentType{}, true, "validation failed: name is a required server component type attribute"},
		{"happy path", &db.ServerComponentType{Name: "Test-Type"}, false, ""},
	}

	for _, tt := range testCases {
		err := s.CreateServerComponentType(tt.ct)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
		}
	}
}

func TestGetServerComponentType(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName      string
		filter        *db.ServerComponentTypeFilter
		expectedUUIDs []uuid.UUID
	}{
		{"happy path - filter doesn't match", &db.ServerComponentTypeFilter{Name: "DoesntExist"}, []uuid.UUID{}},
		{"happy path - filter match", &db.ServerComponentTypeFilter{Name: db.FixtureSCTFins.Name}, []uuid.UUID{db.FixtureSCTFins.ID}},
		{"happy path - no filter", nil, []uuid.UUID{db.FixtureSCTFins.ID}},
	}

	for _, tt := range testCases {
		r, err := s.GetServerComponentTypes(tt.filter, nil)
		assert.NoError(t, err, tt.testName)

		var rIDs []uuid.UUID
		for _, h := range r {
			rIDs = append(rIDs, h.ID)
		}

		assert.ElementsMatch(t, rIDs, tt.expectedUUIDs, tt.testName)
	}
}
