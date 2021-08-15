package gormdb_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/gormdb"
)

func TestGetComponentsByServerUUID(t *testing.T) {
	s := gormdb.DatabaseTest(t)

	var testCases = []struct {
		testName       string
		srvUUID        uuid.UUID
		filter         *gormdb.ServerComponentFilter
		expectedUUIDs  []uuid.UUID
		expectNotFound bool
	}{
		{
			"happy path - no filter",
			gormdb.FixtureServerNemo.ID,
			nil,
			[]uuid.UUID{
				gormdb.FixtureSCNemoLeftFin.ID,
				gormdb.FixtureSCNemoRightFin.ID,
			},
			false,
		},
		{
			"happy path - filter match",
			gormdb.FixtureServerNemo.ID,
			&gormdb.ServerComponentFilter{Serial: "Left"},
			[]uuid.UUID{
				gormdb.FixtureSCNemoLeftFin.ID,
			},
			false,
		},
		{
			"happy path - filter excludes everything",
			gormdb.FixtureServerNemo.ID,
			&gormdb.ServerComponentFilter{Serial: "Doesnt Exist"},
			[]uuid.UUID{},
			false,
		},
		{
			"happy path - unknown server id",
			uuid.New(),
			nil,
			[]uuid.UUID{},
			true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, count, err := s.GetComponentsByServerUUID(tt.srvUUID, tt.filter, nil)
			if tt.expectNotFound {
				assert.Error(t, err)
				assert.ErrorIs(t, err, gormdb.ErrNotFound)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, len(tt.expectedUUIDs), count)

				var rIDs []uuid.UUID
				for _, h := range r {
					rIDs = append(rIDs, h.ID)
				}

				assert.ElementsMatch(t, rIDs, tt.expectedUUIDs)
			}
		})
	}
}
