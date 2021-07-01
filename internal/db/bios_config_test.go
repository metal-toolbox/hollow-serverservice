package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateBIOSConfig(t *testing.T) {
	databaseTest(t)

	var testCases = []struct {
		testName    string
		bc          db.BIOSConfig
		expectError bool
		errorMsg    string
	}{
		{"missing hardware id", db.BIOSConfig{}, true, "validation failed: hardware UUID is a required BIOSConfig attribute"},
		{"hardware id that doesn't exist", db.BIOSConfig{HardwareID: uuid.New()}, true, "hardware UUID not found"},
		{"happy path", db.BIOSConfig{HardwareID: fixtureHardwareNemo.ID}, false, ""},
	}

	for _, tt := range testCases {
		err := db.CreateBIOSConfig(tt.bc)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
		}
	}
}

func TestBIOSConfigList(t *testing.T) {
	databaseTest(t)

	var testCases = []struct {
		testName    string
		searchUUID  uuid.UUID
		expectList  []db.BIOSConfig
		expectError bool
		errorMsg    string
	}{
		{"no results, bad uuid", uuid.New(), []db.BIOSConfig{}, false, ""},
		{"happy path", fixtureBIOSConfigNew.HardwareID, []db.BIOSConfig{fixtureBIOSConfigNew, fixtureBIOSConfig}, false, ""},
	}

	for _, tt := range testCases {
		res, err := db.BIOSConfigList(tt.searchUUID)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			for i, bc := range tt.expectList {
				assert.Equal(t, bc.ID, res[i].ID)
				assert.Equal(t, bc.ConfigValues, res[i].ConfigValues)
			}
		}
	}
}
