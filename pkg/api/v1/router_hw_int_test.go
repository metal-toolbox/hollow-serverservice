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

var testHW = hollow.Hardware{
	UUID:         uuid.New(),
	FacilityCode: "int-test",
	HardwareComponents: []hollow.HardwareComponent{
		{
			Name:   "Intel Xeon 123",
			Model:  "Xeon 123",
			Vendor: "Intel",
			Serial: "987654321",
			Attributes: []hollow.Attributes{
				{
					Namespace: "hollow.integration.test",
					Values:    json.RawMessage([]byte(`{"firmware":1}`)),
				},
			},
			HardwareComponentTypeUUID: db.FixtureHCTFins.ID,
		},
	},
	Attributes: []hollow.Attributes{
		{
			Namespace: "hollow.integration.test",
			Values:    json.RawMessage([]byte(`{"plan_type":"large"}`)),
		},
	},
	VersionedAttributes: []hollow.VersionedAttributes{
		{
			Namespace: "hollow.integration.settings",
			Values:    json.RawMessage([]byte(`{"setting":"enabled"}`)),
		},
	},
}

func TestIntegrationHardwareList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		res, err := s.Client.Hardware.List(ctx, nil)
		if !expectError {
			require.Len(t, res, 3)
		}

		return err
	})

	// These are the same test cases used in db/hardware_test.go
	var testCases = []struct {
		testName      string
		params        *hollow.HardwareListParams
		expectedUUIDs []uuid.UUID
		expectError   bool
		errorMsg      string
	}{
		{
			"search by age less than 7",
			&hollow.HardwareListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:     db.FixtureNamespaceMetadata,
						Keys:          []string{"age"},
						LessThanValue: 7,
					},
				},
			},
			[]uuid.UUID{db.FixtureHardwareNemo.ID},
			false,
			"",
		},
		{
			"search by age greater than 1 and facility code",
			&hollow.HardwareListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:        db.FixtureNamespaceMetadata,
						Keys:             []string{"age"},
						GreaterThanValue: 1,
					},
				},
				FacilityCode: "Dory",
			},
			[]uuid.UUID{db.FixtureHardwareDory.ID},
			false,
			"",
		},
		{
			"search by facility",
			&hollow.HardwareListParams{
				FacilityCode: "Dory",
			},
			[]uuid.UUID{db.FixtureHardwareDory.ID},
			false,
			"",
		},
		{
			"search by type and location from different attributes",
			&hollow.HardwareListParams{
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
			[]uuid.UUID{db.FixtureHardwareDory.ID},
			false,
			"",
		},
		{
			"search by nested tag",
			&hollow.HardwareListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:  db.FixtureNamespaceOtherdata,
						Keys:       []string{"nested", "tag"},
						EqualValue: "finding-nemo",
					},
				},
			},
			[]uuid.UUID{db.FixtureHardwareDory.ID, db.FixtureHardwareNemo.ID, db.FixtureHardwareMarlin.ID},
			false,
			"",
		},
		{
			"search by nested number greater than 1",
			&hollow.HardwareListParams{
				AttributeListParams: []hollow.AttributeListParams{
					{
						Namespace:        db.FixtureNamespaceOtherdata,
						Keys:             []string{"nested", "number"},
						GreaterThanValue: 1,
					},
				},
			},
			[]uuid.UUID{db.FixtureHardwareDory.ID, db.FixtureHardwareMarlin.ID},
			false,
			"",
		},
		{
			"empty search filter",
			nil,
			[]uuid.UUID{db.FixtureHardwareNemo.ID, db.FixtureHardwareDory.ID, db.FixtureHardwareMarlin.ID},
			false,
			"",
		},
		{
			"facility filter that doesn't match",
			&hollow.HardwareListParams{
				FacilityCode: "Neverland",
			},
			[]uuid.UUID{},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes",
			&hollow.HardwareListParams{
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
			[]uuid.UUID{db.FixtureHardwareNemo.ID},
			false,
			"",
		},
		{
			"search by type from attributes and name from versioned attributes, using the not current value, so nothing should return",
			&hollow.HardwareListParams{
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
	}

	for _, tt := range testCases {
		r, err := s.Client.Hardware.List(context.TODO(), tt.params)
		if tt.expectError {
			assert.NoError(t, err)
			continue
		}

		var actual []uuid.UUID

		for _, hw := range r {
			actual = append(actual, hw.UUID)
		}

		assert.ElementsMatch(t, tt.expectedUUIDs, actual)
	}
}

func TestIntegrationHardwareCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		res, err := s.Client.Hardware.Create(ctx, testHW)
		if !expectError {
			assert.NotNil(t, res)
			assert.Equal(t, testHW.UUID.String(), res.String())
		}

		return err
	})

	var testCases = []struct {
		testName string
		hw       *hollow.Hardware
		errorMsg string
	}{
		{
			"fails on a duplicate uuid",
			&hollow.Hardware{
				UUID:         db.FixtureHardwareNemo.ID,
				FacilityCode: "int-test",
			},
			"duplicate key",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := s.Client.Hardware.Create(context.TODO(), *tt.hw)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestIntegrationHardwareDelete(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		return s.Client.Hardware.Delete(ctx, hollow.Hardware{UUID: db.FixtureHardwareNemo.ID})
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
			err := s.Client.Hardware.Delete(context.TODO(), hollow.Hardware{UUID: tt.uuid})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestIntegrationHardwareCreateAndFetchWithAllAttributes(t *testing.T) {
	s := serverTest(t)
	// Attempt to get the testUUID (should return a failure unless somehow we got a collision with fixtures)
	_, err := s.Client.Hardware.Get(context.TODO(), testHW.UUID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")

	_, err = s.Client.Hardware.Create(context.TODO(), testHW)
	assert.NoError(t, err)

	// Get the hardware back and ensure all the things we set are returned
	rHW, err := s.Client.Hardware.Get(context.TODO(), testHW.UUID)
	assert.NoError(t, err)

	assert.Equal(t, rHW.FacilityCode, "int-test")

	assert.Len(t, rHW.HardwareComponents, 1)
	hc := rHW.HardwareComponents[0]
	assert.Equal(t, "Intel Xeon 123", hc.Name)
	assert.Equal(t, "Xeon 123", hc.Model)
	assert.Equal(t, "Intel", hc.Vendor)
	assert.Equal(t, "987654321", hc.Serial)
	assert.Equal(t, db.FixtureHCTFins.ID, hc.HardwareComponentTypeUUID)
	assert.Equal(t, "Fins", hc.HardwareComponentTypeName)

	assert.Len(t, hc.Attributes, 1)
	assert.Equal(t, "hollow.integration.test", hc.Attributes[0].Namespace)
	assert.Equal(t, json.RawMessage([]byte(`{"firmware":1}`)), hc.Attributes[0].Values)

	assert.Len(t, rHW.Attributes, 1)
	assert.Equal(t, "hollow.integration.test", rHW.Attributes[0].Namespace)
	assert.Equal(t, json.RawMessage([]byte(`{"plan_type":"large"}`)), rHW.Attributes[0].Values)

	assert.Len(t, rHW.VersionedAttributes, 1)
	assert.Equal(t, "hollow.integration.settings", rHW.VersionedAttributes[0].Namespace)
	assert.Equal(t, json.RawMessage([]byte(`{"setting":"enabled"}`)), rHW.VersionedAttributes[0].Values)
}

func TestIntegrationHardwareServiceCreateVersionedAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		va := hollow.VersionedAttributes{Namespace: "hollow.integegration.test", Values: json.RawMessage([]byte(`{"test":"integration"}`))}

		res, err := s.Client.Hardware.CreateVersionedAttributes(ctx, db.FixtureHardwareDory.ID, va)
		if !expectError {
			assert.NotNil(t, res)
		}

		return err
	})
}

func TestIntegrationHardwareServiceGetVersionedAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		res, err := s.Client.Hardware.GetVersionedAttributes(ctx, db.FixtureHardwareNemo.ID)
		if !expectError {
			require.Len(t, res, 2)
			assert.Equal(t, db.FixtureNamespaceVersioned, res[0].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"new"}`)), res[0].Values)
			assert.Equal(t, db.FixtureNamespaceVersioned, res[1].Namespace)
			assert.Equal(t, json.RawMessage([]byte(`{"name":"old"}`)), res[1].Values)
		}

		return err
	})
}
