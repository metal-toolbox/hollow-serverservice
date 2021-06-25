package db_test

import (
	"testing"

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
		{"missing name", db.Hardware{}, true, "validation failed: facility is a required hardware attribute"},
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

func TestHardwareList(t *testing.T) {
	databaseTest(t)

	var testCases = []struct {
		testName    string
		expectList  []db.Hardware
		expectError bool
		errorMsg    string
	}{
		{"missing name", []db.Hardware{fixtureHardware}, false, ""},
	}

	for _, tt := range testCases {
		res, err := db.HardwareList()

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			for i, h := range tt.expectList {
				assert.Equal(t, h.ID, res[i].ID)
			}
		}
	}
}
