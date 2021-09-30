package serverservice_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"go.hollow.sh/serverservice/internal/dbtools"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
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
		params        *serverservice.ServerListParams
		expectedUUIDs []string
		expectError   bool
		errorMsg      string
	}{
		{
			"search by age less than 7",
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceMetadata,
						Keys:      []string{"age"},
						Operator:  serverservice.OperatorLessThan,
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
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceMetadata,
						Keys:      []string{"age"},
						Operator:  serverservice.OperatorGreaterThan,
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
			&serverservice.ServerListParams{
				FacilityCode: "Ocean",
			},
			[]string{dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search by type and location from different attributes",
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  serverservice.OperatorEqual,
						Value:     "blue-tang",
					},
					{
						Namespace: dbtools.FixtureNamespaceMetadata,
						Keys:      []string{"location"},
						Operator:  serverservice.OperatorEqual,
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
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"nested", "tag"},
						Operator:  serverservice.OperatorEqual,
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
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"nested", "number"},
						Operator:  serverservice.OperatorGreaterThan,
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
			&serverservice.ServerListParams{
				FacilityCode: "Neverland",
			},
			[]string{},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes",
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  serverservice.OperatorEqual,
						Value:     "clown",
					},
				},
				VersionedAttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  serverservice.OperatorEqual,
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
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  serverservice.OperatorEqual,
						Value:     "clown",
					},
				},
				VersionedAttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  serverservice.OperatorEqual,
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
			&serverservice.ServerListParams{
				ComponentListParams: []serverservice.ServerComponentListParams{
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
			&serverservice.ServerListParams{
				ComponentListParams: []serverservice.ServerComponentListParams{
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
			&serverservice.ServerListParams{
				ComponentListParams: []serverservice.ServerComponentListParams{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  serverservice.OperatorEqual,
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
			&serverservice.ServerListParams{
				ComponentListParams: []serverservice.ServerComponentListParams{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  serverservice.OperatorEqual,
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
			&serverservice.ServerListParams{
				ComponentListParams: []serverservice.ServerComponentListParams{
					{
						Model: "Normal Fin",
						VersionedAttributeListParams: []serverservice.AttributeListParams{
							{
								Namespace: dbtools.FixtureNamespaceVersioned,
								Keys:      []string{"something"},
								Operator:  serverservice.OperatorEqual,
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
			&serverservice.ServerListParams{
				ComponentListParams: []serverservice.ServerComponentListParams{
					{
						Model: "Normal Fin",
						VersionedAttributeListParams: []serverservice.AttributeListParams{
							{
								Namespace: dbtools.FixtureNamespaceVersioned,
								Keys:      []string{"something"},
								Operator:  serverservice.OperatorEqual,
								Value:     "cool",
							},
						},
					},
				},
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  serverservice.OperatorEqual,
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
			&serverservice.ServerListParams{
				ComponentListParams: []serverservice.ServerComponentListParams{
					{
						Model: "A Lucky Fin",
						VersionedAttributeListParams: []serverservice.AttributeListParams{
							{
								Namespace: dbtools.FixtureNamespaceVersioned,
								Keys:      []string{"something"},
								Operator:  serverservice.OperatorEqual,
								Value:     "cool",
							},
						},
					},
				},
				VersionedAttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Keys:      []string{"name"},
						Operator:  serverservice.OperatorEqual,
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
			&serverservice.ServerListParams{
				ComponentListParams: []serverservice.ServerComponentListParams{
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
			&serverservice.ServerListParams{
				VersionedAttributeListParams: []serverservice.AttributeListParams{
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
			&serverservice.ServerListParams{
				VersionedAttributeListParams: []serverservice.AttributeListParams{
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
			&serverservice.ServerListParams{
				VersionedAttributeListParams: []serverservice.AttributeListParams{
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
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
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
			&serverservice.ServerListParams{},
			[]string{dbtools.FixtureNemo.ID, dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID},
			false,
			"",
		},
		{
			"search for server with IncludeDeleted defined",
			&serverservice.ServerListParams{IncludeDeleted: true},
			[]string{dbtools.FixtureNemo.ID, dbtools.FixtureDory.ID, dbtools.FixtureMarlin.ID, dbtools.FixtureChuckles.ID},
			false,
			"",
		},
		{
			"search for devices by attributes that have a type like %lo%",
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  serverservice.OperatorLike,
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
			&serverservice.ServerListParams{
				AttributeListParams: []serverservice.AttributeListParams{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Keys:      []string{"type"},
						Operator:  serverservice.OperatorLike,
						Value:     "%lo",
					},
				},
			},
			[]string{},
			false,
			"",
		},
	}

	boil.DebugMode = true

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, _, err := s.Client.List(context.TODO(), tt.params)
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

	p := &serverservice.ServerListParams{PaginationParams: &serverservice.PaginationParams{Limit: 2, Page: 1}}
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
	s.Client.SetToken(validToken([]string{"read", "write"}))

	r, _, err := s.Client.Get(context.TODO(), uuid.MustParse(dbtools.FixtureNemo.ID))

	assert.NoError(t, err)
	assert.Len(t, r.Attributes, 2)
	assert.Len(t, r.VersionedAttributes, 1)
	assert.JSONEq(t, string(r.VersionedAttributes[0].Data), string(dbtools.FixtureNemoVersionedNew.Data))
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
	s.Client.SetToken(validToken([]string{"read", "write"}))

	// should only return nemo
	r, _, err := s.Client.List(context.TODO(), &serverservice.ServerListParams{FacilityCode: "Sydney"})

	assert.NoError(t, err)
	assert.Len(t, r, 1)
	assert.Len(t, r[0].Attributes, 2)
	assert.Len(t, r[0].VersionedAttributes, 1)
	assert.JSONEq(t, string(r[0].VersionedAttributes[0].Data), string(dbtools.FixtureNemoVersionedNew.Data))
	assert.Len(t, r[0].Components, 2)
}

func TestIntegrationServerCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		testServer := serverservice.Server{
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
		srv      *serverservice.Server
		errorMsg string
	}{
		{
			"fails on a duplicate uuid",
			&serverservice.Server{
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
		_, err := s.Client.Delete(ctx, serverservice.Server{UUID: uuid.MustParse(dbtools.FixtureNemo.ID)})

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
				_, _, err := s.Client.Create(context.TODO(), serverservice.Server{UUID: tt.uuid})
				assert.NoError(t, err)
			}

			_, err := s.Client.Delete(context.TODO(), serverservice.Server{UUID: tt.uuid})
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

		resp, err := s.Client.Update(ctx, uuid.MustParse(dbtools.FixtureDory.ID), serverservice.Server{Name: "The New Dory"})
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

		va := serverservice.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration"}`))}

		resp, err := s.Client.CreateVersionedAttributes(ctx, uuid.New(), va)
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

	va := serverservice.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration"}`))}
	newVA := serverservice.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration", "something":"else"}`))}

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
			require.Len(t, res, 2)
			assert.Equal(t, dbtools.FixtureNamespaceVersioned, res[0].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"new"}`)), res[0].Data)
			assert.Equal(t, dbtools.FixtureNamespaceVersioned, res[1].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"old"}`)), res[1].Data)
		}

		return err
	})
}
