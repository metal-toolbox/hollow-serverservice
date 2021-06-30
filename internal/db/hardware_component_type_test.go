package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateHardwareComponentType(t *testing.T) {
	databaseTest(t)

	var testCases = []struct {
		testName    string
		hct         db.HardwareComponentType
		expectError bool
		errorMsg    string
	}{
		{"missing name", db.HardwareComponentType{}, true, "validation failed: name is a required hardware component type attribute"},
		{"happy path", db.HardwareComponentType{Name: "Test-Type"}, false, ""},
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
