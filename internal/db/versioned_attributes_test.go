package db_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateVersionedAttributes(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		srv         *db.Server
		a           db.VersionedAttributes
		expectError bool
		errorMsg    string
	}{
		{"missing namespace", &db.FixtureServerDory, db.VersionedAttributes{}, true, "validation failed: namespace is a required VersionedAttributes attribute"},
		{"happy path", &db.FixtureServerDory, db.VersionedAttributes{Namespace: "integration.test.createva", Data: datatypes.JSON([]byte(`{"test":"integration"}`))}, false, ""},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			var beforeUpdated time.Time

			if !tt.expectError {
				s, err := s.FindServerByUUID(tt.srv.ID)
				require.NoError(t, err)
				beforeUpdated = s.UpdatedAt

				// sleep to ensure the timestamp will change
				time.Sleep(1 * time.Second)
			}

			va := &tt.a
			err := s.CreateVersionedAttributes(tt.srv, va)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				// ensure the server updated_at timestamp changed
				s, err := s.FindServerByUUID(tt.srv.ID)
				assert.NoError(t, err)
				assert.NotEqual(t, beforeUpdated, s.UpdatedAt)
			}
		})
	}
}

func TestGetVersionedAttributes(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		searchUUID  uuid.UUID
		searchNS    string
		expectList  []db.VersionedAttributes
		expectError bool
		errorMsg    string
	}{
		{"no results, bad uuid", uuid.New(), "namespace", []db.VersionedAttributes{}, false, ""},
		{"no results, bad namespace", db.FixtureServerNemo.ID, "namespace", []db.VersionedAttributes{}, false, ""},
		{"happy path", db.FixtureServerNemo.ID, db.FixtureNamespaceVersioned, []db.VersionedAttributes{db.FixtureVersionedAttributesNew, db.FixtureVersionedAttributesOld}, false, ""},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			res, count, err := s.GetVersionedAttributes(tt.searchUUID, tt.searchNS, nil)

			if tt.expectError {
				assert.Error(t, err, tt.testName)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err, tt.testName)
				assert.EqualValues(t, len(tt.expectList), count)
				for i, bc := range tt.expectList {
					assert.Equal(t, bc.ID, res[i].ID)
					assert.Equal(t, bc.Data, res[i].Data)
				}
			}
		})
	}
}

func TestListVersionedAttributes(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		searchUUID  uuid.UUID
		expectList  []db.VersionedAttributes
		expectError bool
		errorMsg    string
	}{
		{"no results, bad uuid", uuid.New(), []db.VersionedAttributes{}, false, ""},
		{"happy path", db.FixtureServerNemo.ID, []db.VersionedAttributes{db.FixtureVersionedAttributesNew, db.FixtureVersionedAttributesOld}, false, ""},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			res, count, err := s.ListVersionedAttributes(tt.searchUUID, nil)

			if tt.expectError {
				assert.Error(t, err, tt.testName)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err, tt.testName)
				assert.EqualValues(t, len(tt.expectList), count)
				for i, bc := range tt.expectList {
					assert.Equal(t, bc.ID, res[i].ID)
					assert.Equal(t, bc.Data, res[i].Data)
				}
			}
		})
	}
}
