package fleetdbapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

func TestIntegrationServerList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		r, resp, err := s.Client.List(ctx, nil)
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
		params        *fleetdbapi.ServerListParams
		expectedUUIDs []string
		expectError   bool
		errorMsg      string
	}{
		{
			"search by age less than 7",
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceMetadata,
						Keys:      []string{"age"},
						Operator:  fleetdbapi.OperatorLessThan,
						Value:     "7",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by age greater than 11 and facility code",
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceMetadata,
						Keys:      []string{"age"},
						Operator:  fleetdbapi.OperatorGreaterThan,
						Value:     "11",
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
			&fleetdbapi.ServerListParams{
				FacilityCode: "Ocean",
			},
			[]string{dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search by type and location from different attributes",
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "blue-tang",
					},
					{
						Namespace: dbtools.FixtureNamespaceMetadata,
						Keys:      []string{"location"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "East Australian Current",
					},
				},
			},
			[]string{dbtools.FixtureDory.ID},
			false,
			"",
		},
		{
			"search by nested tag",
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"nested", "tag"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "finding-nemo",
					},
				},
			},
			[]string{dbtools.FixtureDory.ID, dbtools.FixtureNemo.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search by nested number greater than 1",
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"nested", "number"},
						Operator:  fleetdbapi.OperatorGreaterThan,
						Value:     "1",
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
			&fleetdbapi.ServerListParams{
				FacilityCode: "Neverland",
			},
			[]string{},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes",
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "clown",
					},
				},
				VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "new",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes, using the not current value, so nothing should return",
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "clown",
					},
				},
				VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "old",
					},
				},
			},
			[]string{},
			false,
			"",
		},
		{
			"search by multiple components of the server",
			&fleetdbapi.ServerListParams{
				ComponentListParams: []fleetdbapi.ServerComponentListParams{
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
			&fleetdbapi.ServerListParams{
				ComponentListParams: []fleetdbapi.ServerComponentListParams{
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
			&fleetdbapi.ServerListParams{
				ComponentListParams: []fleetdbapi.ServerComponentListParams{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "new",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by a single component and versioned attributes of the server",
			&fleetdbapi.ServerListParams{
				ComponentListParams: []fleetdbapi.ServerComponentListParams{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "new",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by a single component and it's versioned attributes of the server",
			&fleetdbapi.ServerListParams{
				ComponentListParams: []fleetdbapi.ServerComponentListParams{
					{
						Model: "Normal Fin",
						VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
							{
								Namespace: dbtools.FixtureNamespaceVersioned,
								Keys:      []string{"something"},
								Operator:  fleetdbapi.OperatorEqual,
								Value:     "cool",
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
			&fleetdbapi.ServerListParams{
				ComponentListParams: []fleetdbapi.ServerComponentListParams{
					{
						Model: "Normal Fin",
						VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
							{
								Namespace: dbtools.FixtureNamespaceVersioned,
								Keys:      []string{"something"},
								Operator:  fleetdbapi.OperatorEqual,
								Value:     "cool",
							},
						},
					},
				},
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "clown",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID},
			false,
			"",
		},
		{
			"search by a component and server versioned attributes of the server",
			&fleetdbapi.ServerListParams{
				ComponentListParams: []fleetdbapi.ServerComponentListParams{
					{
						Model: "A Lucky Fin",
						VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
							{
								Namespace: dbtools.FixtureNamespaceVersioned,
								Keys:      []string{"something"},
								Operator:  fleetdbapi.OperatorEqual,
								Value:     "cool",
							},
						},
					},
				},
				VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  fleetdbapi.OperatorEqual,
						Value:     "old",
					},
				},
			},
			[]string{},
			false,
			"",
		},
		{
			"search by a component slug",
			&fleetdbapi.ServerListParams{
				ComponentListParams: []fleetdbapi.ServerComponentListParams{
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
			&fleetdbapi.ServerListParams{
				VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
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
			&fleetdbapi.ServerListParams{
				VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
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
			&fleetdbapi.ServerListParams{
				VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
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
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceMetadata,
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID, dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search for server without IncludeDeleted defined",
			&fleetdbapi.ServerListParams{},
			[]string{dbtools.FixtureNemo.ID, dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search for server with IncludeDeleted defined",
			&fleetdbapi.ServerListParams{IncludeDeleted: true},
			[]string{dbtools.FixtureNemo.ID, dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID, dbtools.FixtureChuckles.ID},
			false,
			"",
		},
		{
			"search for devices by attributes that have a type like %lo%",
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  fleetdbapi.OperatorLike,
						Value:     "%lo%",
					},
				},
			},
			[]string{dbtools.FixtureNemo.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search for devices by attributes that have a type like %lo",
			&fleetdbapi.ServerListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  fleetdbapi.OperatorLike,
						Value:     "%lo",
					},
				},
			},
			[]string{},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			servers, resp, err := s.Client.List(context.TODO(), tt.params)
			if tt.expectError {
				assert.NoError(t, err)
				return
			}

			var actual []string

			assert.Equal(t, int64(len(servers)), resp.TotalRecordCount)

			for _, srv := range servers {
				actual = append(actual, srv.UUID.String())
			}

			assert.ElementsMatch(t, tt.expectedUUIDs, actual)
		})
	}
}

func TestIntegrationServerListPagination(t *testing.T) {
	s := serverTest(t)
	s.Client.SetToken(validToken(adminScopes))

	p := &fleetdbapi.ServerListParams{PaginationParams: &fleetdbapi.PaginationParams{Limit: 2, Page: 1}}
	r, resp, err := s.Client.List(context.TODO(), p)

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

func TestIntegrationServerGetPreload(t *testing.T) {
	s := serverTest(t)
	s.Client.SetToken(validToken(adminScopes))

	r, _, err := s.Client.Get(context.TODO(), uuid.MustParse(dbtools.FixtureNemo.ID))

	assert.NoError(t, err)
	assert.Len(t, r.Attributes, 2)
	assert.Len(t, r.VersionedAttributes, 2)
	assert.Len(t, r.Components, 2)
	assert.Nil(t, r.DeletedAt, "DeletedAt should be nil for non deleted server")
}

func TestIntegrationServerGetDeleted(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		r, _, err := s.Client.Get(ctx, uuid.MustParse(dbtools.FixtureChuckles.ID))
		if !expectError {
			require.NoError(t, err)
			assert.Equal(t, r.UUID, uuid.MustParse(dbtools.FixtureChuckles.ID), "Expected UUID %s, got %s", dbtools.FixtureChuckles.ID, r.UUID.String())
			assert.Equal(t, r.Name, dbtools.FixtureChuckles.Name.String)
			assert.NotNil(t, r.DeletedAt)
		}

		return err
	})
}

func TestIntegrationServerListPreload(t *testing.T) {
	s := serverTest(t)
	s.Client.SetToken(validToken(adminScopes))

	// should only return nemo
	r, _, err := s.Client.List(context.TODO(), &fleetdbapi.ServerListParams{FacilityCode: "Sydney"})

	assert.NoError(t, err)
	assert.Len(t, r, 1)
	assert.Len(t, r[0].Attributes, 2)
	assert.Len(t, r[0].VersionedAttributes, 2)
	assert.Len(t, r[0].Components, 2)
}

func TestIntegrationServerCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		testServer := fleetdbapi.Server{
			UUID:         uuid.New(),
			Name:         "test-server",
			FacilityCode: "int",
		}

		id, resp, err := s.Client.Create(ctx, testServer)
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
		srv      *fleetdbapi.Server
		errorMsg string
	}{
		{
			"fails on a duplicate uuid",
			&fleetdbapi.Server{
				UUID:         uuid.MustParse(dbtools.FixtureNemo.ID),
				FacilityCode: "int-test",
			},
			"duplicate key",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			_, _, err := s.Client.Create(context.TODO(), *tt.srv)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestIntegrationServerDelete(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)
		_, err := s.Client.Delete(ctx, fleetdbapi.Server{UUID: uuid.MustParse(dbtools.FixtureNemo.ID)})

		return err
	})

	var testCases = []struct {
		testName  string
		uuid      uuid.UUID
		errorMsg  string
		expectErr bool
		create    bool
	}{
		{
			"fails on unknown uuid",
			uuid.New(),
			"resource not found",
			true,
			false,
		},
		{
			"ensure soft deleted server is retrievable",
			uuid.New(),
			"",
			false,
			true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.create {
				_, _, err := s.Client.Create(context.TODO(), fleetdbapi.Server{UUID: tt.uuid})
				assert.NoError(t, err)
			}

			_, err := s.Client.Delete(context.TODO(), fleetdbapi.Server{UUID: tt.uuid})
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)

				server, _, err := s.Client.Get(context.TODO(), tt.uuid)

				assert.NoError(t, err)
				assert.NotNil(t, server)
				assert.NotEqual(t, server.DeletedAt, null.Time{}.Time)
			}
		})
	}
}

func TestIntegrationServerUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		resp, err := s.Client.Update(ctx, uuid.MustParse(dbtools.FixtureDory.ID), fleetdbapi.Server{Name: "The New Dory"})
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

		va := fleetdbapi.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration"}`))}

		resp, err := s.Client.CreateVersionedAttributes(ctx, uuid.New(), va)
		if !expectError {
			assert.Equal(t, va.Namespace, resp.Slug)
		}

		return err
	})
}

func TestIntegrationServerServiceCreateVersionedAttributesIncrementCounter(t *testing.T) {
	s := serverTest(t)
	s.Client.SetToken(validToken(adminScopes))

	u := uuid.New()
	ctx := context.TODO()

	va := fleetdbapi.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration"}`))}
	newVA := fleetdbapi.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration", "something":"else"}`))}

	_, err := s.Client.CreateVersionedAttributes(ctx, u, va)
	require.NoError(t, err)

	// Ensure we only have one versioned attribute now
	r, _, err := s.Client.GetVersionedAttributes(ctx, u, "hollow.integegration.test")
	require.NoError(t, err)
	assert.Len(t, r, 1)

	// Create with the same data again. This should just increase the counter, not create a new one.
	_, err = s.Client.CreateVersionedAttributes(ctx, u, va)
	require.NoError(t, err)

	// Ensure we still have only one versioned attribute and that the counter increased
	r, _, err = s.Client.GetVersionedAttributes(ctx, u, "hollow.integegration.test")
	require.NoError(t, err)
	assert.Len(t, r, 1)
	assert.Equal(t, 1, r[0].Tally)

	// Create with different data and ensure a new one is created
	_, err = s.Client.CreateVersionedAttributes(ctx, u, newVA)
	require.NoError(t, err)

	// Ensure we still have only one versioned attribute and that the counter increased
	r, _, err = s.Client.GetVersionedAttributes(ctx, u, "hollow.integegration.test")
	require.NoError(t, err)
	assert.Len(t, r, 2)
	assert.Equal(t, 0, r[0].Tally)
	assert.Equal(t, 1, r[1].Tally)
}

func TestIntegrationServerServiceListVersionedAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		res, _, err := s.Client.ListVersionedAttributes(ctx, uuid.MustParse(dbtools.FixtureNemo.ID))
		if !expectError {
			require.Len(t, res, 3)
			assert.Equal(t, dbtools.FixtureNamespaceVersioned, res[0].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"new"}`)), res[0].Data)
			assert.Equal(t, dbtools.FixtureNamespaceVersioned, res[1].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"old"}`)), res[1].Data)
		}

		return err
	})
}
