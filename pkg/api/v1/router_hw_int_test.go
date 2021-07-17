package hollow_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestIntegrationHardwareList(t *testing.T) {
	s := serverTest(t)

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

	var testCases = []struct {
		testName    string
		hw          *hollow.Hardware
		expectError bool
		errorMsg    string
	}{
		{
			"happy path",
			&hollow.Hardware{
				FacilityCode: "int-test",
			},
			false,
			"",
		},
		{
			"fails on a duplicate uuid",
			&hollow.Hardware{
				UUID:         db.FixtureHardwareNemo.ID,
				FacilityCode: "int-test",
			},
			true,
			"duplicate key",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, err := s.Client.Hardware.Create(context.TODO(), *tt.hw)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.NotEqual(t, uuid.Nil.String(), r.String())
			}
		})
	}
}

func TestIntegrationHardwareDelete(t *testing.T) {
	s := serverTest(t)

	var testCases = []struct {
		testName    string
		uuid        uuid.UUID
		expectError bool
		errorMsg    string
	}{
		{
			"happy path",
			db.FixtureHardwareNemo.ID,
			false,
			"",
		},
		{
			"fails on unknown uuid",
			uuid.New(),
			true,
			"resource not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			err := s.Client.Hardware.Delete(context.TODO(), hollow.Hardware{UUID: tt.uuid})
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIntegrationHardwareCreateAndFetchWithAllAttributes(t *testing.T) {
	s := serverTest(t)
	testUUID := uuid.New()

	// Attempt to get the testUUID (should return a failure unless somehow we got a collision with fixtures)
	_, err := s.Client.Hardware.Get(context.TODO(), testUUID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")

	hw := &hollow.Hardware{
		UUID:         testUUID,
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
	}

	_, err = s.Client.Hardware.Create(context.TODO(), *hw)
	assert.NoError(t, err)

	// Get the hardware back and ensure all the things we set are returned
	rHW, err := s.Client.Hardware.Get(context.TODO(), testUUID)
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
}

func TestIntegrationHardwareServiceCreateVersionedAttributes(t *testing.T) {
	s := serverTest(t)
	hwUUID := db.FixtureHardwareDory.ID

	var testCases = []struct {
		testName    string
		va          hollow.VersionedAttributes
		expectError bool
		errorMsg    string
	}{
		{
			"happy path",
			hollow.VersionedAttributes{
				Namespace: "hollow.integration.test",
				Values:    json.RawMessage([]byte(`{"integration":true}`)),
			},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, err := s.Client.Hardware.CreateVersionedAttributes(context.TODO(), hwUUID, tt.va)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.NotEqual(t, uuid.Nil.String(), r.String())
			}
		})
	}
}
