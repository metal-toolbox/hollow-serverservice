package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
)

func TestGetComponentsByServerUUID(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName       string
		srvUUID        uuid.UUID
		filter         *db.ServerComponentFilter
		expectedUUIDs  []uuid.UUID
		expectNotFound bool
	}{
		{
			"happy path - no filter",
			db.FixtureServerNemo.ID,
			nil,
			[]uuid.UUID{
				db.FixtureSCNemoLeftFin.ID,
				db.FixtureSCNemoRightFin.ID,
			},
			false,
		},
		{
			"happy path - filter match",
			db.FixtureServerNemo.ID,
			&db.ServerComponentFilter{Serial: "Left"},
			[]uuid.UUID{
				db.FixtureSCNemoLeftFin.ID,
			},
			false,
		},
		{
			"happy path - filter excludes everything",
			db.FixtureServerNemo.ID,
			&db.ServerComponentFilter{Serial: "Doesnt Exist"},
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
				assert.ErrorIs(t, err, db.ErrNotFound)
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
