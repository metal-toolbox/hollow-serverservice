package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateHardwareComponentType(t *testing.T) {
	db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		hct         *db.HardwareComponentType
		expectError bool
		errorMsg    string
	}{
		{"missing name", &db.HardwareComponentType{}, true, "validation failed: name is a required hardware component type attribute"},
		{"happy path", &db.HardwareComponentType{Name: "Test-Type"}, false, ""},
	}

	for _, tt := range testCases {
		err := db.CreateHardwareComponentType(tt.hct)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
		}
	}
}

func TestGetHardwareComponentType(t *testing.T) {
	db.DatabaseTest(t)

	var testCases = []struct {
		testName      string
		filter        *db.HardwareComponentTypeFilter
		expectedUUIDs []uuid.UUID
		expectError   bool
		errorMsg      string
	}{
		{"happy path - filter doesn't match", &db.HardwareComponentTypeFilter{Name: "DoesntExist"}, []uuid.UUID{}, false, ""},
		{"happy path - filter match", &db.HardwareComponentTypeFilter{Name: db.FixtureHCTFins.Name}, []uuid.UUID{db.FixtureHCTFins.ID}, false, ""},
		{"happy path - no filter", nil, []uuid.UUID{db.FixtureHCTFins.ID}, false, ""},
	}

	for _, tt := range testCases {
		r, err := db.GetHardwareComponentTypes(tt.filter)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)

			var rIDs []uuid.UUID
			for _, h := range r {
				rIDs = append(rIDs, h.ID)
			}

			assert.ElementsMatch(t, rIDs, tt.expectedUUIDs, tt.testName)
		}
	}
}
