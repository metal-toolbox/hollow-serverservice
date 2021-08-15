package gormdb_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/gormdb"
)

func TestCreateServerComponentType(t *testing.T) {
	s := gormdb.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		ct          *gormdb.ServerComponentType
		expectError bool
		errorMsg    string
	}{
		{"missing name", &gormdb.ServerComponentType{}, true, "validation failed: name is a required server component type attribute"},
		{"happy path", &gormdb.ServerComponentType{Name: "Test-Type"}, false, ""},
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
	s := gormdb.DatabaseTest(t)

	var testCases = []struct {
		testName      string
		filter        *gormdb.ServerComponentTypeFilter
		expectedUUIDs []uuid.UUID
	}{
		{"happy path - filter doesn't match", &gormdb.ServerComponentTypeFilter{Name: "DoesntExist"}, []uuid.UUID{}},
		{"happy path - filter match", &gormdb.ServerComponentTypeFilter{Name: gormdb.FixtureSCTFins.Name}, []uuid.UUID{gormdb.FixtureSCTFins.ID}},
		{"happy path - no filter", nil, []uuid.UUID{gormdb.FixtureSCTFins.ID}},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, count, err := s.GetServerComponentTypes(tt.filter, nil)
			assert.NoError(t, err, tt.testName)
			assert.EqualValues(t, len(tt.expectedUUIDs), count)

			var rIDs []uuid.UUID
			for _, h := range r {
				rIDs = append(rIDs, h.ID)
			}

			assert.ElementsMatch(t, rIDs, tt.expectedUUIDs, tt.testName)
		})
	}
}
