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
		{"missing name", db.BIOSConfig{}, true, "validation failed: hardware UUID is a required BIOSConfig attribute"},
		{"happy path", db.BIOSConfig{HardwareUUID: uuid.New()}, false, ""},
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
		{"happy path", fixtureBIOSConfigNew.HardwareUUID, []db.BIOSConfig{fixtureBIOSConfigNew, fixtureBIOSConfigOld}, false, ""},
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
