package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateHardware(t *testing.T) {
	db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		hw          db.Hardware
		expectError bool
		errorMsg    string
	}{
		// {"missing name", db.Hardware{}, true, "validation failed: facility is a required hardware attribute"},
		{"happy path", db.Hardware{FacilityCode: "TEST1"}, false, ""},
	}

	for _, tt := range testCases {
		err := db.CreateHardware(&tt.hw)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
		}
	}
}

func TestFindOrCreateHardwareByUUID(t *testing.T) {
	db.DatabaseTest(t)

	var testCases = []struct {
		testName   string
		searchUUID uuid.UUID

		expectError bool
		errorMsg    string
	}{
		{
			"happy path - existing hardware",
			db.FixtureHardwareDory.ID,
			false,
			"",
		},
		{
			"happy path - new hardware",
			uuid.New(),
			false,
			"",
		},
	}

	for _, tt := range testCases {
		res, err := db.FindOrCreateHardwareByUUID(tt.searchUUID)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, res, tt.testName)
			assert.NotNil(t, res.CreatedAt, tt.testName)
		}
	}
}

func TestFindHardwareByUUID(t *testing.T) {
	db.DatabaseTest(t)

	var testCases = []struct {
		testName   string
		searchUUID uuid.UUID

		expectError bool
		errorMsg    string
	}{
		{
			"happy path - existing hardware",
			db.FixtureHardwareDory.ID,
			false,
			"",
		},
		{
			"happy path - hardware not found",
			uuid.New(),
			true,
			"something not found",
		},
	}

	for _, tt := range testCases {
		res, err := db.FindHardwareByUUID(tt.searchUUID)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Errorf(t, err, tt.errorMsg, tt.testName)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, res, tt.testName)
			assert.NotNil(t, res.CreatedAt, tt.testName)
		}
	}
}

func TestGetHardware(t *testing.T) {
	db.DatabaseTest(t)

	var testCases = []struct {
		testName      string
		filter        *db.HardwareFilter
		expectedUUIDs []uuid.UUID
		expectError   bool
		errorMsg      string
	}{
		{
			"search by age less than 7",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:     db.FixtureNamespaceMetadata,
						Keys:          []string{"age"},
						LessThanValue: 7,
					},
				},
			},
			[]uuid.UUID{db.FixtureHardwareNemo.ID},
			false,
			"",
		},
		{
			"search by age greater than 1 and facility code",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:        db.FixtureNamespaceMetadata,
						Keys:             []string{"age"},
						GreaterThanValue: 1,
					},
				},
				FacilityCode: "Dory",
			},
			[]uuid.UUID{db.FixtureHardwareDory.ID},
			false,
			"",
		},
		{
			"search by facility",
			&db.HardwareFilter{
				FacilityCode: "Dory",
			},
			[]uuid.UUID{db.FixtureHardwareDory.ID},
			false,
			"",
		},
		{
			"search by type and location from different attributes",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "blue-tang",
					},
					{
						Namespace:  db.FixtureNamespaceMetadata,
						Keys:       []string{"location"},
						EqualValue: "East Austalian Current",
					},
				},
			},
			[]uuid.UUID{db.FixtureHardwareDory.ID},
			false,
			"",
		},
		{
			"search by nested tag",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"nested", "tag"},
						EqualValue: "finding-nemo",
					},
				},
			},
			[]uuid.UUID{db.FixtureHardwareDory.ID, db.FixtureHardwareNemo.ID, db.FixtureHardwareMarlin.ID},
			false,
			"",
		},
		{
			"search by nested number greater than 1",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:        db.FixtureNamespaceOtherdata,
						Keys:             []string{"nested", "number"},
						GreaterThanValue: 1,
					},
				},
			},
			[]uuid.UUID{db.FixtureHardwareDory.ID, db.FixtureHardwareMarlin.ID},
			false,
			"",
		},
		{
			"empty search filter",
			nil,
			[]uuid.UUID{db.FixtureHardwareNemo.ID, db.FixtureHardwareDory.ID, db.FixtureHardwareMarlin.ID},
			false,
			"",
		},
		{
			"facility filter that doesn't match",
			&db.HardwareFilter{
				FacilityCode: "Neverland",
			},
			[]uuid.UUID{},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]uuid.UUID{db.FixtureHardwareNemo.ID},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes, using the not current value, so nothing should return",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "old",
					},
				},
			},
			[]uuid.UUID{},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		r, err := db.GetHardware(tt.filter)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg, tt.testName)
		} else {
			assert.NoError(t, err)

			var rIDs []uuid.UUID
			for _, h := range r {
				rIDs = append(rIDs, h.ID)
				// Ensure preload works. All Fixture data has 2 hardware components
				assert.Len(t, h.HardwareComponents, 2, tt.testName)
			}

			assert.ElementsMatch(t, rIDs, tt.expectedUUIDs, tt.testName)
		}
	}
}
