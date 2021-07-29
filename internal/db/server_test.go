package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
)

func TestCreateServer(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName    string
		srv         db.Server
		expectError bool
		errorMsg    string
	}{
		{"happy path", db.Server{FacilityCode: "TEST1"}, false, ""},
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
	s := db.DatabaseTest(t)

	err := s.DeleteServer(&db.FixtureServerNemo)
	assert.NoError(t, err)
}

func TestFindServerByUUID(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName   string
		searchUUID uuid.UUID

		expectError bool
		errorMsg    string
	}{
		{
			"happy path - existing server",
			db.FixtureServerDory.ID,
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
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName   string
		searchUUID uuid.UUID

		expectError bool
		errorMsg    string
	}{
		{
			"happy path - existing server",
			db.FixtureServerDory.ID,
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
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName      string
		filter        *db.ServerFilter
		expectedUUIDs []uuid.UUID
		expectError   bool
		errorMsg      string
	}{
		{
			"search by age less than 7",
			&db.ServerFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:     db.FixtureNamespaceMetadata,
						Keys:          []string{"age"},
						LessThanValue: 7,
					},
				},
			},
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by age greater than 11 and facility code",
			&db.ServerFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:        db.FixtureNamespaceMetadata,
						Keys:             []string{"age"},
						GreaterThanValue: 11,
					},
				},
				FacilityCode: "Ocean",
			},
			[]uuid.UUID{db.FixtureServerDory.ID},
			false,
			"",
		},
		{
			"search by facility",
			&db.ServerFilter{
				FacilityCode: "Ocean",
			},
			[]uuid.UUID{db.FixtureServerDory.ID, db.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"search by type and location from different attributes",
			&db.ServerFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "blue-tang",
					},
					{
						Namespace:  db.FixtureNamespaceMetadata,
						Keys:       []string{"location"},
						EqualValue: "East Austalian Current",
					},
				},
			},
			[]uuid.UUID{db.FixtureServerDory.ID},
			false,
			"",
		},
		{
			"search by nested tag",
			&db.ServerFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"nested", "tag"},
						EqualValue: "finding-nemo",
					},
				},
			},
			[]uuid.UUID{db.FixtureServerDory.ID, db.FixtureServerNemo.ID, db.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"search by nested number greater than 1",
			&db.ServerFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:        db.FixtureNamespaceOtherdata,
						Keys:             []string{"nested", "number"},
						GreaterThanValue: 1,
					},
				},
			},
			[]uuid.UUID{db.FixtureServerDory.ID, db.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"empty search filter",
			nil,
			[]uuid.UUID{db.FixtureServerNemo.ID, db.FixtureServerDory.ID, db.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"facility filter that doesn't match",
			&db.ServerFilter{
				FacilityCode: "Neverland",
			},
			[]uuid.UUID{},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes",
			&db.ServerFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes, using the not current value, so nothing should return",
			&db.ServerFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceVersioned,
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
			&db.ServerFilter{
				ComponentFilters: []db.ServerComponentFilter{
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
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"ensure both components have to match when searching by multiple components of the server",
			&db.ServerFilter{
				ComponentFilters: []db.ServerComponentFilter{
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
			&db.ServerFilter{
				ComponentFilters: []db.ServerComponentFilter{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by a single component and versioned attributes of the server",
			&db.ServerFilter{
				ComponentFilters: []db.ServerComponentFilter{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by a single component and it's versioned attributes of the server",
			&db.ServerFilter{
				ComponentFilters: []db.ServerComponentFilter{
					{
						Model: "A Lucky Fin",
						VersionedAttributesFilters: []db.AttributesFilter{
							{
								Namespace:  db.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
			},
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by a component and server versioned attributes of the server",
			&db.ServerFilter{
				ComponentFilters: []db.ServerComponentFilter{
					{
						Model: "A Lucky Fin",
						VersionedAttributesFilters: []db.AttributesFilter{
							{
								Namespace:  db.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search by a component and server versioned attributes of the server",
			&db.ServerFilter{
				ComponentFilters: []db.ServerComponentFilter{
					{
						Model: "A Lucky Fin",
						VersionedAttributesFilters: []db.AttributesFilter{
							{
								Namespace:  db.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace:  db.FixtureNamespaceVersioned,
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
			&db.ServerFilter{
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace: db.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
					},
				},
			},
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search for devices with a versioned attributes in a namespace with key that doesn't exists",
			&db.ServerFilter{
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace: db.FixtureNamespaceVersioned,
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
			&db.ServerFilter{
				VersionedAttributesFilters: []db.AttributesFilter{
					{
						Namespace: db.FixtureNamespaceVersioned,
					},
				},
			},
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"search for devices that have attributes in a namespace - no filters",
			&db.ServerFilter{
				AttributesFilters: []db.AttributesFilter{
					{
						Namespace: db.FixtureNamespaceMetadata,
					},
				},
			},
			[]uuid.UUID{db.FixtureServerNemo.ID, db.FixtureServerDory.ID, db.FixtureServerMarlin.ID},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, err := s.GetServers(tt.filter, nil)

			if tt.expectError {
				assert.Error(t, err, tt.testName)
				assert.Contains(t, err.Error(), tt.errorMsg, tt.testName)
			} else {
				assert.NoError(t, err)

				var rIDs []uuid.UUID
				for _, h := range r {
					rIDs = append(rIDs, h.ID)
					// Ensure preload works. All Fixture data has 2 server components and 2 attributes
					assert.Len(t, h.ServerComponents, 2, tt.testName)
					assert.Len(t, h.Attributes, 2, tt.testName)
					// Nemo has two versioned attributes but only the most recent in a namespace should preload
					if h.ID == db.FixtureServerNemo.ID {
						assert.Len(t, h.VersionedAttributes, 1, tt.testName)
					}
				}

				assert.ElementsMatch(t, rIDs, tt.expectedUUIDs, tt.testName)
			}
		})
	}
}

func TestGetServerPagination(t *testing.T) {
	s := db.DatabaseTest(t)

	var testCases = []struct {
		testName      string
		pager         db.Pagination
		expectedUUIDs []uuid.UUID
		expectError   bool
		errorMsg      string
	}{
		{
			"limit 1 page 1",
			db.Pagination{
				Limit: 1,
				Page:  1,
				Sort:  "created_at DESC",
			},
			[]uuid.UUID{db.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"limit 1 page 2",
			db.Pagination{
				Limit: 1,
				Page:  2,
				Sort:  "created_at DESC",
			},
			[]uuid.UUID{db.FixtureServerDory.ID},
			false,
			"",
		},
		{
			"limit 1 page 3",
			db.Pagination{
				Limit: 1,
				Page:  3,
				Sort:  "created_at DESC",
			},
			[]uuid.UUID{db.FixtureServerNemo.ID},
			false,
			"",
		},
		{
			"limit 1 page 4",
			db.Pagination{
				Limit: 1,
				Page:  4,
				Sort:  "created_at DESC",
			},
			[]uuid.UUID{},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		r, err := s.GetServers(nil, &tt.pager)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg, tt.testName)
		} else {
			assert.NoError(t, err)

			var rIDs []uuid.UUID
			for _, h := range r {
				rIDs = append(rIDs, h.ID)
			}

			assert.ElementsMatch(t, rIDs, tt.expectedUUIDs, tt.testName)
		}
	}
}
