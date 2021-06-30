package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateHardware(t *testing.T) {
	databaseTest(t)

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
		err := db.CreateHardware(tt.hw)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
		}
	}
}

func TestFindOrCreateHardwareByUUID(t *testing.T) {
	databaseTest(t)

	var testCases = []struct {
		testName   string
		searchUUID uuid.UUID

		expectError bool
		errorMsg    string
	}{
		{
			"happy path - existing hardware",
			fixtureHardwareDory.ID,
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

func TestGetHardwareWithFilter(t *testing.T) {
	databaseTest(t)

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
						Namespace:     fixtureNamespaceMetadata,
						Keys:          []string{"age"},
						LessThanValue: 7,
					},
				},
			},
			[]uuid.UUID{fixtureHardwareNemo.ID},
			false,
			"",
		},
		{
			"search by age greater than 1 and facility code",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:        fixtureNamespaceMetadata,
						Keys:             []string{"age"},
						GreaterThanValue: 1,
					},
				},
				FacilityCode: "Dory",
			},
			[]uuid.UUID{fixtureHardwareDory.ID},
			false,
			"",
		},
		{
			"search by facility",
			&db.HardwareFilter{
				FacilityCode: "Dory",
			},
			[]uuid.UUID{fixtureHardwareDory.ID},
			false,
			"",
		},
		{
			"search by type and location from different attributes",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  fixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "blue-tang",
					},
					{
						Namespace:  fixtureNamespaceMetadata,
						Keys:       []string{"location"},
						EqualValue: "East Austalian Current",
					},
				},
			},
			[]uuid.UUID{fixtureHardwareDory.ID},
			false,
			"",
		},
		{
			"search by nested tag",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  fixtureNamespaceOtherdata,
						Keys:       []string{"nested", "tag"},
						EqualValue: "finding-nemo",
					},
				},
			},
			[]uuid.UUID{fixtureHardwareDory.ID, fixtureHardwareNemo.ID, fixtureHardwareMarlin.ID},
			false,
			"",
		},
		{
			"search by nested number greater than 1",
			&db.HardwareFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:        fixtureNamespaceOtherdata,
						Keys:             []string{"nested", "number"},
						GreaterThanValue: 1,
					},
				},
			},
			[]uuid.UUID{fixtureHardwareDory.ID, fixtureHardwareMarlin.ID},
			false,
			"",
		},
		{
			"empty search filter",
			nil,
			[]uuid.UUID{fixtureHardwareNemo.ID, fixtureHardwareDory.ID, fixtureHardwareMarlin.ID},
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
			}

			assert.ElementsMatch(t, rIDs, tt.expectedUUIDs, tt.testName)
		}
	}
}
