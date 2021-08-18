package hollow_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"go.metalkube.net/hollow/internal/dbtools"
	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

// var testServer = hollow.Server{
// 	UUID:         uuid.New(),
// 	FacilityCode: "int-test",
// 	Components: []hollow.ServerComponent{
// 		{
// 			Name:   "Intel Xeon 123",
// 			Model:  "Xeon 123",
// 			Vendor: "Intel",
// 			Serial: "987654321",
// 			Attributes: []hollow.Attributes{
// 				{
// 					Namespace: "hollow.integration.test",
// 					Data:      json.RawMessage([]byte(`{"firmware":1}`)),
// 				},
// 			},
// 			ComponentTypeID: gormdb.FixtureSCTFins.Slug,
// 		},
// 	},
// 	Attributes: []hollow.Attributes{
// 		{
// 			Namespace: "hollow.integration.test",
// 			Data:      json.RawMessage([]byte(`{"plan_type":"large"}`)),
// 		},
// 	},
// 	VersionedAttributes: []hollow.VersionedAttributes{
// 		{
// 			Namespace: "hollow.integration.settings",
// 			Data:      json.RawMessage([]byte(`{"setting":"enabled"}`)),
// 		},
// 	},
// }

func TestIntegrationServerList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		r, resp, err := s.Client.Server.List(ctx, nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, r, 3)

			assert.EqualValues(t, 3, resp.PageCount)
			assert.EqualValues(t, 1, resp.TotalPages)
			assert.EqualValues(t, 3, resp.TotalRecordCount)
			// We returned everything, so we shouldnt have a next page info
			assert.Nil(t, resp.Links.Next)
			assert.Nil(t, resp.Links.Previous)
		}

		return err
	})

	// These are the same test cases used in db/server_test.go
	var testCases = []struct {
		testName      string
		params        *hollow.ServerListParams
		expectedUUIDs []string
		expectError   bool
		errorMsg      string
	}{
		{
			"search by age less than 7",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:     dbtools.FixtureNamespaceMetadata,
						Keys:          []string{"age"},
						LessThanValue: 7,
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by age greater than 11 and facility code",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:        dbtools.FixtureNamespaceMetadata,
						Keys:             []string{"age"},
						GreaterThanValue: 11,
					},
				},
				FacilityCode: "Ocean",
			},
			[]string{dbtools.FixtureDory.ID},
			false,
			"",
		},
		{
			"search by facility",
			&hollow.ServerListParams{
				FacilityCode: "Ocean",
			},
			[]string{dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search by type and location from different attributes",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "blue-tang",
					},
					{
						Namespace:  dbtools.FixtureNamespaceMetadata,
						Keys:       []string{"location"},
						EqualValue: "East Australian Current",
					},
				},
			},
			[]string{dbtools.FixtureDory.ID},
			false,
			"",
		},
		{
			"search by nested tag",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceOtherdata,
						Keys:       []string{"nested", "tag"},
						EqualValue: "finding-nemo",
					},
				},
			},
			[]string{dbtools.FixtureDory.ID, dbtools.FixtureNemo.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search by nested number greater than 1",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:        dbtools.FixtureNamespaceOtherdata,
						Keys:             []string{"nested", "number"},
						GreaterThanValue: 1,
					},
				},
			},
			[]string{dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"empty search filter",
			nil,
			[]string{dbtools.FixtureNemo.ID, dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"facility filter that doesn't match",
			&hollow.ServerListParams{
				FacilityCode: "Neverland",
			},
			[]string{},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes, using the not current value, so nothing should return",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "old",
					},
				},
			},
			[]string{},
			false,
			"",
		},
		{
			"search by multiple components of the server",
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
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
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"ensure both components have to match when searching by multiple components of the server",
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
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
			[]string{},
			false,
			"",
		},
		{
			"search by a single component and versioned attributes of the server",
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by a single component and versioned attributes of the server",
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "new",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by a single component and it's versioned attributes of the server",
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						Model: "Normal Fin",
						VersionedAttributeListParams: []hollow.AttributeListParams{
							{
								Namespace:  dbtools.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by a component and server attributes of the server",
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						Model: "Normal Fin",
						VersionedAttributeListParams: []hollow.AttributeListParams{
							{
								Namespace:  dbtools.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by a component and server versioned attributes of the server",
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						Model: "A Lucky Fin",
						VersionedAttributeListParams: []hollow.AttributeListParams{
							{
								Namespace:  dbtools.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  dbtools.FixtureNamespaceVersioned,
						Keys:       []string{"name"},
						EqualValue: "old",
					},
				},
			},
			[]string{},
			false,
			"",
		},
		{
			"search by a component slug",
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						ServerComponentType: dbtools.FixtureFinType.Slug,
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID, dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search for devices with a versioned attributes in a namespace with key that exists",
			&hollow.ServerListParams{
				VersionedAttributeListParams: []hollow.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search for devices with a versioned attributes in a namespace with key that doesn't exists",
			&hollow.ServerListParams{
				VersionedAttributeListParams: []hollow.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"doesntExist"},
					},
				},
			},
			[]string{},
			false,
			"",
		},
		{
			"search for devices that have versioned attributes in a namespace - no filters",
			&hollow.ServerListParams{
				VersionedAttributeListParams: []hollow.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search for devices that have attributes in a namespace - no filters",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceMetadata,
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID, dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
	}

	boil.DebugMode = true

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, _, err := s.Client.Server.List(context.TODO(), tt.params)
			if tt.expectError {
				assert.NoError(t, err)
				return
			}

			var actual []string

			for _, srv := range r {
				actual = append(actual, srv.UUID.String())
			}

			assert.ElementsMatch(t, tt.expectedUUIDs, actual)
		})
	}

	boil.DebugMode = false
}

func TestIntegrationServerListPagination(t *testing.T) {
	s := serverTest(t)
	s.Client.SetToken(validToken([]string{"read", "write"}))

	p := &hollow.ServerListParams{PaginationParams: &hollow.PaginationParams{Limit: 2, Page: 1}}
	r, resp, err := s.Client.Server.List(context.TODO(), p)

	assert.NoError(t, err)
	assert.Len(t, r, 2)
	assert.Equal(t, dbtools.FixtureServers[2].ID, r[0].UUID.String())
	assert.Equal(t, dbtools.FixtureServers[1].ID, r[1].UUID.String())

	assert.EqualValues(t, 2, resp.PageCount)
	assert.EqualValues(t, 2, resp.TotalPages)
	assert.EqualValues(t, 3, resp.TotalRecordCount)
	// Since we have a next page let's make sure all the links are set
	assert.NotNil(t, resp.Links.Next)
	assert.Nil(t, resp.Links.Previous)
	assert.True(t, resp.HasNextPage())

	//
	// Get the next page and verify the results
	//
	resp, err = s.Client.NextPage(context.TODO(), *resp, &r)

	assert.NoError(t, err)
	assert.Len(t, r, 1)
	assert.Equal(t, dbtools.FixtureServers[0].ID, r[0].UUID.String())

	assert.EqualValues(t, 1, resp.PageCount)

	// we should have followed the cursor so first/previous/next/last links shouldn't be set
	// but there is another page so we should have a next cursor link. Total counts are not includes
	// cursor pages
	assert.EqualValues(t, 2, resp.TotalPages)
	assert.EqualValues(t, 3, resp.TotalRecordCount)
	assert.NotNil(t, resp.Links.First)
	assert.NotNil(t, resp.Links.Previous)
	assert.Nil(t, resp.Links.Next)
	assert.NotNil(t, resp.Links.Last)
	assert.False(t, resp.HasNextPage())
}

func TestIntegrationServerCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		testServer := hollow.Server{
			UUID:         uuid.New(),
			Name:         "test-server",
			FacilityCode: "int",
		}

		id, resp, err := s.Client.Server.Create(ctx, testServer)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, id)
			assert.Equal(t, testServer.UUID.String(), id.String())
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/servers/%s", id), resp.Links.Self.Href)
		}

		return err
	})

	var testCases = []struct {
		testName string
		srv      *hollow.Server
		errorMsg string
	}{
		{
			"fails on a duplicate uuid",
			&hollow.Server{
				UUID:         uuid.MustParse(dbtools.FixtureNemo.ID),
				FacilityCode: "int-test",
			},
			"duplicate key",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			_, _, err := s.Client.Server.Create(context.TODO(), *tt.srv)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestIntegrationServerDelete(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)
		_, err := s.Client.Server.Delete(ctx, hollow.Server{UUID: uuid.MustParse(dbtools.FixtureNemo.ID)})

		return err
	})

	var testCases = []struct {
		testName string
		uuid     uuid.UUID
		errorMsg string
	}{
		{
			"fails on unknown uuid",
			uuid.New(),
			"resource not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := s.Client.Server.Delete(context.TODO(), hollow.Server{UUID: tt.uuid})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestIntegrationServerUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		resp, err := s.Client.Server.Update(ctx, uuid.MustParse(dbtools.FixtureDory.ID), hollow.Server{Name: "The New Dory"})
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/servers/%s", dbtools.FixtureDory.ID), resp.Links.Self.Href)
		}

		return err
	})
}

func TestIntegrationServerServiceCreateVersionedAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		va := hollow.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration"}`))}

		resp, err := s.Client.Server.CreateVersionedAttributes(ctx, uuid.New(), va)
		if !expectError {
			assert.Equal(t, va.Namespace, resp.Slug)
		}

		return err
	})
}

func TestIntegrationServerServiceCreateVersionedAttributesIncrementCounter(t *testing.T) {
	s := serverTest(t)
	s.Client.SetToken(validToken([]string{"read", "write"}))

	u := uuid.New()
	ctx := context.TODO()

	va := hollow.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration"}`))}
	newVA := hollow.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration", "something":"else"}`))}

	_, err := s.Client.Server.CreateVersionedAttributes(ctx, u, va)
	require.NoError(t, err)

	// Ensure we only have one versioned attribute now
	r, _, err := s.Client.Server.GetVersionedAttributes(ctx, u, "hollow.integegration.test")
	require.NoError(t, err)
	assert.Len(t, r, 1)

	// Create with the same data again. This should just increase the counter, not create a new one.
	_, err = s.Client.Server.CreateVersionedAttributes(ctx, u, va)
	require.NoError(t, err)

	// Ensure we still have only one versioned attribute and that the counter increased
	r, _, err = s.Client.Server.GetVersionedAttributes(ctx, u, "hollow.integegration.test")
	require.NoError(t, err)
	assert.Len(t, r, 1)
	assert.Equal(t, 1, r[0].Tally)

	// Create with different data and ensure a new one is created
	_, err = s.Client.Server.CreateVersionedAttributes(ctx, u, newVA)
	require.NoError(t, err)

	// Ensure we still have only one versioned attribute and that the counter increased
	r, _, err = s.Client.Server.GetVersionedAttributes(ctx, u, "hollow.integegration.test")
	require.NoError(t, err)
	assert.Len(t, r, 2)
	assert.Equal(t, 0, r[0].Tally)
	assert.Equal(t, 1, r[1].Tally)
}

func TestIntegrationServerServiceListVersionedAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		res, _, err := s.Client.Server.ListVersionedAttributes(ctx, uuid.MustParse(dbtools.FixtureNemo.ID))
		if !expectError {
			require.Len(t, res, 2)
			assert.Equal(t, dbtools.FixtureNamespaceVersioned, res[0].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"new"}`)), res[0].Data)
			assert.Equal(t, dbtools.FixtureNamespaceVersioned, res[1].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"old"}`)), res[1].Data)
		}

		return err
	})
}
