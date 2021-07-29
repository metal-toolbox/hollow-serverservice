package hollow_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.metalkube.net/hollow/internal/db"
	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

var testServer = hollow.Server{
	UUID:         uuid.New(),
	FacilityCode: "int-test",
	Components: []hollow.ServerComponent{
		{
			Name:   "Intel Xeon 123",
			Model:  "Xeon 123",
			Vendor: "Intel",
			Serial: "987654321",
			Attributes: []hollow.Attributes{
				{
					Namespace: "hollow.integration.test",
					Data:      json.RawMessage([]byte(`{"firmware":1}`)),
				},
			},
			ComponentTypeUUID: db.FixtureSCTFins.ID,
		},
	},
	Attributes: []hollow.Attributes{
		{
			Namespace: "hollow.integration.test",
			Data:      json.RawMessage([]byte(`{"plan_type":"large"}`)),
		},
	},
	VersionedAttributes: []hollow.VersionedAttributes{
		{
			Namespace: "hollow.integration.settings",
			Data:      json.RawMessage([]byte(`{"setting":"enabled"}`)),
		},
	},
}

func TestIntegrationServerList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		res, err := s.Client.Server.List(ctx, nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, res, 3)
		}

		return err
	})

	// These are the same test cases used in db/server_test.go
	var testCases = []struct {
		testName      string
		params        *hollow.ServerListParams
		expectedUUIDs []uuid.UUID
		expectError   bool
		errorMsg      string
	}{
		{
			"search by age less than 7",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				FacilityCode: "Ocean",
			},
			[]uuid.UUID{db.FixtureServerDory.ID, db.FixtureServerMarlin.ID},
			false,
			"",
		},
		{
			"search by type and location from different attributes",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				FacilityCode: "Neverland",
			},
			[]uuid.UUID{},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes",
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"type"},
						EqualValue: "clown",
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
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
			[]uuid.UUID{db.FixtureServerNemo.ID},
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
			[]uuid.UUID{},
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
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						Model:  "A Lucky Fin",
						Serial: "Right",
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						Model: "A Lucky Fin",
						VersionedAttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						Model: "A Lucky Fin",
						VersionedAttributeListParams: []hollow.AttributeListParams{
							{
								Namespace:  db.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				ComponentListParams: []hollow.ServerComponentListParams{
					{
						Model: "A Lucky Fin",
						VersionedAttributeListParams: []hollow.AttributeListParams{
							{
								Namespace:  db.FixtureNamespaceVersioned,
								Keys:       []string{"something"},
								EqualValue: "cool",
							},
						},
					},
				},
				VersionedAttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				VersionedAttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				VersionedAttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				VersionedAttributeListParams: []hollow.AttributeListParams{
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
			&hollow.ServerListParams{
				AttributeListParams: []hollow.AttributeListParams{
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
			r, err := s.Client.Server.List(context.TODO(), tt.params)
			if tt.expectError {
				assert.NoError(t, err)
				return
			}

			var actual []uuid.UUID

			for _, srv := range r {
				actual = append(actual, srv.UUID)
			}

			assert.ElementsMatch(t, tt.expectedUUIDs, actual)
		})
	}
}

func TestIntegrationServerCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		res, err := s.Client.Server.Create(ctx, testServer)
		if !expectError {
			assert.NotNil(t, res)
			assert.Equal(t, testServer.UUID.String(), res.String())
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
				UUID:         db.FixtureServerNemo.ID,
				FacilityCode: "int-test",
			},
			"duplicate key",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := s.Client.Server.Create(context.TODO(), *tt.srv)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestIntegrationServerDelete(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		return s.Client.Server.Delete(ctx, hollow.Server{UUID: db.FixtureServerNemo.ID})
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
			err := s.Client.Server.Delete(context.TODO(), hollow.Server{UUID: tt.uuid})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestIntegrationServerCreateAndFetchWithAllAttributes(t *testing.T) {
	s := serverTest(t)
	s.Client.SetToken(validToken([]string{"read", "write"}))

	// Attempt to get the testUUID (should return a failure unless somehow we got a collision with fixtures)
	_, err := s.Client.Server.Get(context.TODO(), testServer.UUID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")

	_, err = s.Client.Server.Create(context.TODO(), testServer)
	assert.NoError(t, err)

	// Get the server back and ensure all the things we set are returned
	r, err := s.Client.Server.Get(context.TODO(), testServer.UUID)
	assert.NoError(t, err)

	assert.Equal(t, r.FacilityCode, "int-test")

	assert.Len(t, r.Components, 1)
	hc := r.Components[0]
	assert.Equal(t, "Intel Xeon 123", hc.Name)
	assert.Equal(t, "Xeon 123", hc.Model)
	assert.Equal(t, "Intel", hc.Vendor)
	assert.Equal(t, "987654321", hc.Serial)
	assert.Equal(t, db.FixtureSCTFins.ID, hc.ComponentTypeUUID)
	assert.Equal(t, "Fins", hc.ComponentTypeName)

	assert.Len(t, hc.Attributes, 1)
	assert.Equal(t, "hollow.integration.test", hc.Attributes[0].Namespace)
	assert.Equal(t, json.RawMessage([]byte(`{"firmware":1}`)), hc.Attributes[0].Data)

	assert.Len(t, r.Attributes, 1)
	assert.Equal(t, "hollow.integration.test", r.Attributes[0].Namespace)
	assert.Equal(t, json.RawMessage([]byte(`{"plan_type":"large"}`)), r.Attributes[0].Data)

	assert.Len(t, r.VersionedAttributes, 1)
	assert.Equal(t, "hollow.integration.settings", r.VersionedAttributes[0].Namespace)
	assert.Equal(t, json.RawMessage([]byte(`{"setting":"enabled"}`)), r.VersionedAttributes[0].Data)
}

func TestIntegrationServerServiceCreateVersionedAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		va := hollow.VersionedAttributes{Namespace: "hollow.integegration.test", Data: json.RawMessage([]byte(`{"test":"integration"}`))}

		res, err := s.Client.Server.CreateVersionedAttributes(ctx, uuid.New(), va)
		if !expectError {
			assert.NotNil(t, res)
		}

		return err
	})
}

func TestIntegrationServerServiceGetVersionedAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		res, err := s.Client.Server.GetVersionedAttributes(ctx, db.FixtureServerNemo.ID)
		if !expectError {
			require.Len(t, res, 2)
			assert.Equal(t, db.FixtureNamespaceVersioned, res[0].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"new"}`)), res[0].Data)
			assert.Equal(t, db.FixtureNamespaceVersioned, res[1].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"old"}`)), res[1].Data)
		}

		return err
	})
}
