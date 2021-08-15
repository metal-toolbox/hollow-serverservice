package gormdb_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/gormdb"
)

func TestCreateServer(t *testing.T) {
	s := gormdb.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		srv         gormdb.Server
		expectError bool
		errorMsg    string
	}{
		{"happy path", gormdb.Server{FacilityCode: "TEST1"}, false, ""},
	}

	for _, tt := range testCases {
		err := s.CreateServer(&tt.srv)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
		}
	}
}

func TestDeleteServer(t *testing.T) {
	s := gormdb.DatabaseTest(t)

	err := s.DeleteServer(&gormdb.FixtureServerNemo)
	assert.NoError(t, err)
}

func TestFindServerByUUID(t *testing.T) {
	s := gormdb.DatabaseTest(t)

	var testCases = []struct {
		testName   string
		searchUUID uuid.UUID

		expectError bool
		errorMsg    string
	}{
		{
			"happy path - existing server",
			gormdb.FixtureServerDory.ID,
			false,
			"",
		},
		{
			"happy path - server not found",
			uuid.New(),
			true,
			"something not found",
		},
	}

	for _, tt := range testCases {
		res, err := s.FindServerByUUID(tt.searchUUID)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Errorf(t, err, tt.errorMsg, tt.testName)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, res, tt.testName)
			assert.NotNil(t, res.CreatedAt, tt.testName)
			assert.Equal(t, tt.searchUUID.String(), res.ID.String())
		}
	}
}

func TestFindOrCreateServerByUUID(t *testing.T) {
	s := gormdb.DatabaseTest(t)

	var testCases = []struct {
		testName   string
		searchUUID uuid.UUID

		expectError bool
		errorMsg    string
	}{
		{
			"happy path - existing server",
			gormdb.FixtureServerDory.ID,
			false,
			"",
		},
		{
			"happy path - server not found, new one created",
			uuid.New(),
			false,
			"",
		},
	}

	for _, tt := range testCases {
		res, err := s.FindOrCreateServerByUUID(tt.searchUUID)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Errorf(t, err, tt.errorMsg, tt.testName)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, res, tt.testName)
			assert.NotNil(t, res.CreatedAt, tt.testName)
			assert.Equal(t, tt.searchUUID.String(), res.ID.String())
		}
	}
}

func TestGetServer(t *testing.T) {
	s := gormdb.DatabaseTest(t)

	var testCases = []struct {
		testName      string
		filter        *gormdb.ServerFilter
		expectedUUIDs []uuid.UUID
		expectError   bool
		errorMsg      string
	}{
		{
			"search by age less than 7",
			&gormdb.ServerFilter{
				AttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:     gormdb.FixtureNamespaceMetadata,
						Keys:          []string{"age"},
						LessThanValue: 7,
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by age greater than 11 and facility code",
			&gormdb.ServerFilter{
				AttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:        gormdb.FixtureNamespaceMetadata,
						Keys:             []string{"age"},
						GreaterThanValue: 11,
					},
				},
				FacilityCode: "Ocean",
			},
			[]uuid.UUID{gormdb.FixtureServerDory.ID},
			false,
			"",
		},
		{
			"search by facility",
			&gormdb.ServerFilter{
				FacilityCode: "Ocean",
			},
			[]uuid.UUID{gormdb.FixtureServerDory.ID, gormdb.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"search by type and location from different attributes",
			&gormdb.ServerFilter{
				AttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "blue-tang",
					},
					{
						Namespace:  gormdb.FixtureNamespaceMetadata,
						Keys:       []string{"location"},
						EqualValue: "East Austalian Current",
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerDory.ID},
			false,
			"",
		},
		{
			"search by nested tag",
			&gormdb.ServerFilter{
				AttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceOtherdata,
						Keys:       []string{"nested", "tag"},
						EqualValue: "finding-nemo",
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerDory.ID, gormdb.FixtureServerNemo.ID, gormdb.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"search by nested number greater than 1",
			&gormdb.ServerFilter{
				AttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:        gormdb.FixtureNamespaceOtherdata,
						Keys:             []string{"nested", "number"},
						GreaterThanValue: 1,
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerDory.ID, gormdb.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"empty search filter",
			nil,
			[]uuid.UUID{gormdb.FixtureServerNemo.ID, gormdb.FixtureServerDory.ID, gormdb.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"facility filter that doesn't match",
			&gormdb.ServerFilter{
				FacilityCode: "Neverland",
			},
			[]uuid.UUID{},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes",
			&gormdb.ServerFilter{
				AttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes, using the not current value, so nothing should return",
			&gormdb.ServerFilter{
				AttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "old",
					},
				},
			},
			[]uuid.UUID{},
			false,
			"",
		},
		{
			"search by multiple components of the server",
			&gormdb.ServerFilter{
				ComponentFilters: []gormdb.ServerComponentFilter{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
					{
						Model:  "Normal Fin",
						Serial: "Left",
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"ensure both components have to match when searching by multiple components of the server",
			&gormdb.ServerFilter{
				ComponentFilters: []gormdb.ServerComponentFilter{
					{
						Name:   "My Lucky Fin",
						Vendor: "Barracuda",
						Model:  "A Lucky Fin",
						Serial: "Left",
					},
					{
						Model:  "Normal Fin",
						Serial: "Left",
					},
				},
			},
			[]uuid.UUID{},
			false,
			"",
		},
		{
			"search by a single component and versioned attributes of the server",
			&gormdb.ServerFilter{
				ComponentFilters: []gormdb.ServerComponentFilter{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by a single component and versioned attributes of the server",
			&gormdb.ServerFilter{
				ComponentFilters: []gormdb.ServerComponentFilter{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by a single component and it's versioned attributes of the server",
			&gormdb.ServerFilter{
				ComponentFilters: []gormdb.ServerComponentFilter{
					{
						Model: "A Lucky Fin",
						VersionedAttributesFilters: []gormdb.AttributesFilter{
							{
								Namespace:  gormdb.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by a component and server versioned attributes of the server",
			&gormdb.ServerFilter{
				ComponentFilters: []gormdb.ServerComponentFilter{
					{
						Model: "A Lucky Fin",
						VersionedAttributesFilters: []gormdb.AttributesFilter{
							{
								Namespace:  gormdb.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
				VersionedAttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by a component and server versioned attributes of the server",
			&gormdb.ServerFilter{
				ComponentFilters: []gormdb.ServerComponentFilter{
					{
						Model: "A Lucky Fin",
						VersionedAttributesFilters: []gormdb.AttributesFilter{
							{
								Namespace:  gormdb.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
				VersionedAttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace:  gormdb.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "old",
					},
				},
			},
			[]uuid.UUID{},
			false,
			"",
		},
		{
			"search for devices with a versioned attributes in a namespace with key that exists",
			&gormdb.ServerFilter{
				VersionedAttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace: gormdb.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search for devices with a versioned attributes in a namespace with key that doesn't exists",
			&gormdb.ServerFilter{
				VersionedAttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace: gormdb.FixtureNamespaceVersioned,
						Keys:      []string{"doesntExist"},
					},
				},
			},
			[]uuid.UUID{},
			false,
			"",
		},
		{
			"search for devices that have versioned attributes in a namespace - no filters",
			&gormdb.ServerFilter{
				VersionedAttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace: gormdb.FixtureNamespaceVersioned,
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search for devices that have attributes in a namespace - no filters",
			&gormdb.ServerFilter{
				AttributesFilters: []gormdb.AttributesFilter{
					{
						Namespace: gormdb.FixtureNamespaceMetadata,
					},
				},
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID, gormdb.FixtureServerDory.ID, gormdb.FixtureServerMarlin.ID},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, count, err := s.GetServers(tt.filter, nil)

			if tt.expectError {
				assert.Error(t, err, tt.testName)
				assert.Contains(t, err.Error(), tt.errorMsg, tt.testName)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, len(tt.expectedUUIDs), count)

				var rIDs []uuid.UUID
				for _, h := range r {
					rIDs = append(rIDs, h.ID)
					// Ensure preload works. All Fixture data has 2 server components and 2 attributes
					assert.Len(t, h.ServerComponents, 2, tt.testName)
					assert.Len(t, h.Attributes, 2, tt.testName)
					// Nemo has two versioned attributes but only the most recent in a namespace should preload
					if h.ID == gormdb.FixtureServerNemo.ID {
						assert.Len(t, h.VersionedAttributes, 1, tt.testName)
					}
				}

				assert.ElementsMatch(t, rIDs, tt.expectedUUIDs, tt.testName)
			}
		})
	}
}

func TestGetServerPagination(t *testing.T) {
	s := gormdb.DatabaseTest(t)

	var testCases = []struct {
		testName      string
		pager         gormdb.Pagination
		expectedUUIDs []uuid.UUID
		expectError   bool
		errorMsg      string
	}{
		{
			"limit 1 page 1",
			gormdb.Pagination{
				Limit: 1,
				Page:  1,
			},
			[]uuid.UUID{gormdb.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"limit 1 page 2",
			gormdb.Pagination{
				Limit: 1,
				Page:  2,
			},
			[]uuid.UUID{gormdb.FixtureServerDory.ID},
			false,
			"",
		},
		{
			"limit 1 page 3",
			gormdb.Pagination{
				Limit: 1,
				Page:  3,
			},
			[]uuid.UUID{gormdb.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"limit 1 page 4",
			gormdb.Pagination{
				Limit: 1,
				Page:  4,
			},
			[]uuid.UUID{},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, count, err := s.GetServers(nil, &tt.pager)

			if tt.expectError {
				assert.Error(t, err, tt.testName)
				assert.Contains(t, err.Error(), tt.errorMsg, tt.testName)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, 3, count)

				var rIDs []uuid.UUID
				for _, h := range r {
					rIDs = append(rIDs, h.ID)
				}

				assert.ElementsMatch(t, rIDs, tt.expectedUUIDs, tt.testName)
			}
		})
	}
}

func TestUpdateServer(t *testing.T) {
	s := gormdb.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		srvUUID     uuid.UUID
		newValues   gormdb.Server
		expectError bool
		errorMsg    string
	}{
		{
			"happy path - existing server",
			gormdb.FixtureServerDory.ID,
			gormdb.Server{Name: "Not Dory", FacilityCode: "Somewhere"},
			false,
			"",
		},
		{
			"no changes - existing server",
			gormdb.FixtureServerDory.ID,
			gormdb.Server{Name: gormdb.FixtureServerDory.Name, FacilityCode: gormdb.FixtureServerDory.FacilityCode},
			false,
			"",
		},
		{
			"server not found",
			uuid.New(),
			gormdb.Server{},
			true,
			"something not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			err := s.UpdateServer(tt.srvUUID, tt.newValues)

			if tt.expectError {
				assert.Error(t, err)
				assert.Errorf(t, err, tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
