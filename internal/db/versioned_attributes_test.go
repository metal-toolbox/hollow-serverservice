package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateVersionedAttributes(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		a           *db.VersionedAttributes
		expectError bool
		errorMsg    string
	}{
		{"missing namespace", &db.VersionedAttributes{}, true, "validation failed: namespace is a required VersionedAttribute attribute"},
		{"happy path", &db.VersionedAttributes{EntityType: "hardware", EntityID: db.FixtureHardwareNemo.ID, Namespace: db.FixtureNamespaceVersioned}, false, ""},
	}

	for _, tt := range testCases {
		err := s.CreateVersionedAttributes(tt.a)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
		}
	}
}

func TestGetVersionedAttributes(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		searchUUID  uuid.UUID
		expectList  []db.VersionedAttributes
		expectError bool
		errorMsg    string
	}{
		{"no results, bad uuid", uuid.New(), []db.VersionedAttributes{}, false, ""},
		{"happy path", db.FixtureVersionedAttributesNew.EntityID, []db.VersionedAttributes{db.FixtureVersionedAttributesNew, db.FixtureVersionedAttributesOld}, false, ""},
	}

	for _, tt := range testCases {
		res, err := s.GetVersionedAttributes(tt.searchUUID)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			for i, bc := range tt.expectList {
				assert.Equal(t, bc.ID, res[i].ID)
				assert.Equal(t, bc.Values, res[i].Values)
			}
		}
	}
}
