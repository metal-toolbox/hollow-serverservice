package fleetdbapi_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

// r640FirmwareFixtureUUIDs returns firmware  uuids based on the firmware hardware model attribute
func r640FirmwareFixtureUUIDs(t *testing.T, firmware []fleetdbapi.ComponentFirmwareVersion) []string {
	t.Helper()

	ids := []string{}

	for idx, f := range firmware {
		if slices.Contains(f.Model, "R640") {
			ids = append(ids, firmware[idx].UUID.String())
		}
	}

	return ids
}

func TestIntegrationServerComponentFirmwareSetCreate(t *testing.T) {
	s := serverTest(t)

	var firmwareSetID uuid.UUID

	firmwareSetID, err := uuid.Parse(dbtools.FixtureFirmwareSetR640.ID)
	if err != nil {
		t.Fatal(err)
	}

	var r640FirmwareIDs []string

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		var testFirmwareSet fleetdbapi.ComponentFirmwareSetRequest

		if !expectError {
			// 2. retrieve component firmware fixture data for test
			firmware, _, err := s.Client.GetServerComponentFirmwareSet(context.Background(), firmwareSetID)
			if err != nil {
				t.Fatal(err)
			}

			assert.NotNil(t, firmware)

			r640FirmwareIDs = r640FirmwareFixtureUUIDs(t, firmware.ComponentFirmware)

			// expect two fixture firmware objects to be returned
			assert.Equal(t, 2, len(r640FirmwareIDs))

			testFirmwareSet = fleetdbapi.ComponentFirmwareSetRequest{
				Name:                   "test-firmware-set",
				ComponentFirmwareUUIDs: r640FirmwareIDs,
			}
		}

		id, resp, err := s.Client.CreateServerComponentFirmwareSet(ctx, testFirmwareSet)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, id)
			assert.Equal(t, "resource created", resp.Message)
			assert.NotNil(t, resp.Links.Self)
		}

		return err
	})

	var testCases = []struct {
		testName           string
		firmwareSetPayload *fleetdbapi.ComponentFirmwareSetRequest
		expectedError      bool
		expectedResponse   string
		errorMsg           string
	}{
		{
			"Name field required",
			&fleetdbapi.ComponentFirmwareSetRequest{},
			true,
			"400",
			"required attribute not set: Name",
		},
		{
			"component firmware UUIDs required",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar",
			},
			true,
			"400",
			"expected one or more firmware UUIDs, got none",
		},
		{
			"valid UUIDs for the firmware ID expected",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar",
				ComponentFirmwareUUIDs: []string{
					r640FirmwareIDs[0],
					"d825bbeb-20fb-452e-9fe4-invalid",
				},
			},
			true,
			"400",
			"invalid firmware UUID",
		},
		{
			"duplicate firmware UUIDs are not accepted",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar",
				ComponentFirmwareUUIDs: []string{
					r640FirmwareIDs[0],
					r640FirmwareIDs[0],
				},
			},
			true,
			"400",
			"A firmware set can only reference unique firmware versions",
		},
		{
			"non-existing firmware UUID are not accepted",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar",
				ComponentFirmwareUUIDs: []string{
					"d825bbeb-20fb-452e-9fe4-cdedacb2ca1f",
				},
			},
			true,
			"400",
			"firmware object with given UUID does not exist",
		},
		{
			"firmware set added referencing firmware UUIDs",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name:                   "foobar",
				ComponentFirmwareUUIDs: r640FirmwareIDs,
				Attributes: []fleetdbapi.Attributes{
					{
						Namespace: "sh.hollow.firmware_set.metadata",
						Data:      json.RawMessage(`{"created by": "foobar"}`),
					},
				},
			},
			false,
			"200",
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			fwUUID, resp, err := s.Client.CreateServerComponentFirmwareSet(context.TODO(), *tt.firmwareSetPayload)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Contains(t, err.Error(), tt.expectedResponse)
				return
			}

			assert.NotNil(t, resp)
			assert.Equal(t, "resource created", resp.Message)
			assert.NotEqual(t, uuid.Nil, fwUUID)
		})
	}
}

func TestIntegrationServerComponentFirmwareSetUpdate(t *testing.T) {
	s := serverTest(t)

	var firmwareSetID uuid.UUID

	firmwareSetID, err := uuid.Parse(dbtools.FixtureFirmwareSetR640.ID)
	if err != nil {
		t.Fatal(err)
	}

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		var err error

		_, err = s.Client.UpdateComponentFirmwareSetRequest(ctx, firmwareSetID, fleetdbapi.ComponentFirmwareSetRequest{})
		if !expectError {
			return nil
		}

		return err
	})

	// retrieve component firmware fixture data for test
	firmware, _, err := s.Client.GetServerComponentFirmwareSet(context.Background(), firmwareSetID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, firmware)

	r640FirmwareIDs := r640FirmwareFixtureUUIDs(t, firmware.ComponentFirmware)

	var testCases = []struct {
		testName                    string
		firmwareSetPayload          *fleetdbapi.ComponentFirmwareSetRequest
		expectFirmwareSetAttributes []fleetdbapi.Attributes
		expectedFirmwareCount       int
		expectedError               bool
		expectedResponse            string
		errorMsg                    string
	}{
		{
			"component firmware set UUID required",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar",
			},
			nil,
			0,
			true,
			"400",
			"expected a valid firmware set ID, got none",
		},
		{
			"valid UUIDs for the firmware ID expected",
			&fleetdbapi.ComponentFirmwareSetRequest{
				ID:   firmwareSetID,
				Name: "foobar",
				ComponentFirmwareUUIDs: []string{
					"d825bbeb-20fb-452e-9fe4-invalid",
				},
			},
			nil,
			0,
			true,
			"400",
			"invalid firmware UUID",
		},
		{
			"duplicate firmware UUIDs are not accepted",
			&fleetdbapi.ComponentFirmwareSetRequest{
				ID:   firmwareSetID,
				Name: "foobar",
				ComponentFirmwareUUIDs: []string{
					r640FirmwareIDs[0],
					r640FirmwareIDs[0],
				},
			},
			nil,
			0,
			true,
			"400",
			"exists in firmware set",
		},
		{
			"non-existing firmware UUID are not accepted",
			&fleetdbapi.ComponentFirmwareSetRequest{
				ID:   firmwareSetID,
				Name: "foobar",
				ComponentFirmwareUUIDs: []string{
					"d825bbeb-20fb-452e-9fe4-cdedacb2ca1f",
				},
			},
			nil,
			0,
			true,
			"400",
			"firmware object with given UUID does not exist",
		},
		{
			"update an existing firmware set - update name, referenced firmware",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name:                   "foobar-updated",
				ID:                     firmwareSetID,
				ComponentFirmwareUUIDs: []string{dbtools.FixtureDellR640CPLD.ID},
			},
			nil,
			3,
			false,
			"200",
			"",
		},
		{
			"update an existing firmware set - update labels",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar-updated",
				ID:   firmwareSetID,
				Attributes: []fleetdbapi.Attributes{
					{
						Namespace: "sh.hollow.firmware_set.labels",
						Data:      json.RawMessage(`{"created by": "foobar"}`),
					},
				},
			},
			[]fleetdbapi.Attributes{
				{
					Namespace: "sh.hollow.firmware_set.labels",
					Data:      json.RawMessage(`{"created by": "foobar"}`),
				},
			},
			3,
			false,
			"200",
			"",
		},
		{
			"update a firmware set with no attributes data does not overwrite existing attributes",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar-updated",
				ID:   firmwareSetID,
				Attributes: []fleetdbapi.Attributes{
					{
						Namespace: "sh.hollow.firmware_set.labels",
					},
				},
			},
			[]fleetdbapi.Attributes{
				{
					Namespace: "sh.hollow.firmware_set.labels",
					Data:      json.RawMessage(`{"created by": "foobar"}`),
				},
			},
			3,
			false,
			"200",
			"",
		},
		{
			"update a firmware set with new attributes updates existing attributes",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar-updated",
				ID:   firmwareSetID,
				Attributes: []fleetdbapi.Attributes{
					{
						Namespace: "sh.hollow.firmware_set.labels",
						Data:      json.RawMessage(`{"updated by": "foo"}`),
					},
				},
			},
			[]fleetdbapi.Attributes{
				{
					Namespace: "sh.hollow.firmware_set.labels",
					Data:      json.RawMessage(`{"updated by": "foo"}`),
				},
			},
			3,
			false,
			"200",
			"",
		},
		{
			"update an existing firmware set with empty attributes does not purge existing attributes",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar-updated-again",
				ID:   firmwareSetID,
			},
			[]fleetdbapi.Attributes{
				{
					Namespace: "sh.hollow.firmware_set.labels",
					Data:      json.RawMessage(`{"updated by": "foo"}`),
				},
			},
			3,
			false,
			"200",
			"",
		},
		{
			"update a firmware set attributes to be empty",
			&fleetdbapi.ComponentFirmwareSetRequest{
				Name: "foobar-updated",
				ID:   firmwareSetID,
				Attributes: []fleetdbapi.Attributes{
					{
						Namespace: "sh.hollow.firmware_set.labels",
						Data:      json.RawMessage(`{}`),
					},
				},
			},
			[]fleetdbapi.Attributes{
				{
					Namespace: "sh.hollow.firmware_set.labels",
					Data:      json.RawMessage(`{}`),
				},
			},
			3,
			false,
			"200",
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			resp, err := s.Client.UpdateComponentFirmwareSetRequest(context.TODO(), firmwareSetID, *tt.firmwareSetPayload)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Contains(t, err.Error(), tt.expectedResponse)
				return
			}

			assert.NotNil(t, resp)
			assert.Equal(t, "resource updated", resp.Message)

			// query firmware set and assert attributes are updated
			got, _, err := s.Client.GetServerComponentFirmwareSet(context.TODO(), firmwareSetID)
			if err != nil {
				t.Fatal(err)
			}

			assert.NotNil(t, got)
			assert.Equal(t, got.UUID.String(), firmwareSetID.String())
			assert.Equal(t, tt.firmwareSetPayload.Name, got.Name)
			assert.Equal(t, tt.expectedFirmwareCount, len(got.ComponentFirmware))

			// assert firmware set attributes
			if len(tt.firmwareSetPayload.Attributes) > 0 {
				assert.Equal(t, 1, len(got.Attributes))
				assert.Equal(t, tt.firmwareSetPayload.Attributes[0].Namespace, got.Attributes[0].Namespace)
			}

			if tt.expectFirmwareSetAttributes != nil {
				assert.NotNil(t, got.Attributes, "Attributes nil")
				assertAttributesEqual(t, tt.expectFirmwareSetAttributes[0].Data, got.Attributes[0].Data)
			}

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Contains(t, err.Error(), tt.expectedResponse)
				return
			}
		})
	}
}

func TestIntegrationServerComponentFirmwareSetGet(t *testing.T) {
	s := serverTest(t)

	var firmwareSetID uuid.UUID

	firmwareSetID, err := uuid.Parse(dbtools.FixtureFirmwareSetR640.ID)
	if err != nil {
		t.Fatal(err)
	}

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		var err error

		_, _, err = s.Client.GetServerComponentFirmwareSet(ctx, firmwareSetID)
		if !expectError {
			require.NoError(t, err)
		}

		return err
	})

	var testCases = []struct {
		testName         string
		firmwareSetID    uuid.UUID
		expectedError    bool
		expectedResponse string
		errorMsg         string
	}{

		{
			"component firmware set UUID required",
			uuid.Nil,
			true,
			"400",
			"expected a firmware set UUID, got none",
		},
		{
			"404 returned for unknown firmware set UUID",
			uuid.New(),
			true,
			"404",
			"resource not found",
		},
		{
			"get an existing firmware set",
			firmwareSetID,
			false,
			"200",
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			got, resp, err := s.Client.GetServerComponentFirmwareSet(context.TODO(), tt.firmwareSetID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Contains(t, err.Error(), tt.expectedResponse)
				return
			}

			assert.NotNil(t, resp)
			assert.Equal(t, got.UUID, tt.firmwareSetID)

			// assert firmware set attributes
			assert.Equal(t, 1, len(got.Attributes))
			assert.Equal(t, dbtools.FixtureFirmwareSetR640Attribute.Namespace, got.Attributes[0].Namespace)
			assertAttributesContains(t, got.Attributes, []byte(dbtools.FixtureFirmwareSetR640Attribute.Data))
			assertAttributesContains(t, got.Attributes, []byte(dbtools.FixtureFirmwareSetX11DPHTAttribute.Data))

			// assert component firmware
			assert.Equal(t, 2, len(got.ComponentFirmware))
			assert.Equal(t, "r640", got.Name)
		})
	}
}

func assertAttributesContains(t *testing.T, attrs []fleetdbapi.Attributes, a []byte) bool {
	for _, attr := range attrs {
		if assertAttributesEqual(t, a, attr.Data) {
			return true
		}
	}

	return false
}

func assertAttributesEqual(t *testing.T, a, b []byte) bool {
	t.Helper()

	// unmarshal fixture attribute data
	aData := map[string]string{}
	if err := json.Unmarshal(a, &aData); err != nil {
		t.Error(err)
	}

	// unmarshal got attribute data
	bData := map[string]string{}
	if err := json.Unmarshal(b, &bData); err != nil {
		t.Error(err)
	}

	return maps.Equal(aData, bData)
}

func TestIntegrationServerComponentFirmwareSetList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		_, _, err := s.Client.ListServerComponentFirmwareSet(ctx, nil)
		if !expectError {
			require.NoError(t, err)
		}

		return err
	})

	testCases := []struct {
		testName                     string
		params                       *fleetdbapi.ComponentFirmwareSetListParams
		expectedFirmwareSetAttribute []*models.AttributesFirmwareSet
		expectedFirmwareModels       []string
		expectedTotalRecordCount     int
		expectedPage                 int
		expectedError                bool
		errorMsg                     string
	}{

		{
			"list firmware set by name - r640",
			&fleetdbapi.ComponentFirmwareSetListParams{Name: "r640"},
			[]*models.AttributesFirmwareSet{dbtools.FixtureFirmwareSetR640Attribute},
			[]string{"R640"},
			1,
			1,
			false,
			"",
		},
		{
			"list firmware set by name - r6515",
			&fleetdbapi.ComponentFirmwareSetListParams{Name: "r6515"},
			[]*models.AttributesFirmwareSet{dbtools.FixtureFirmwareSetR6515Attribute},
			[]string{"R6515"},
			1,
			1,
			false,
			"",
		},
		{
			"list with pagination Limit, Offset",
			&fleetdbapi.ComponentFirmwareSetListParams{
				Pagination: &fleetdbapi.PaginationParams{
					Limit: 1,
					Page:  2,
				},
			},
			nil,
			nil,
			3,
			2,
			false,
			"",
		},
		{
			"list firmware set by attribute params",
			&fleetdbapi.ComponentFirmwareSetListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: "sh.hollow.firmware_set.labels",
						Keys:      []string{"vendor"},
						Operator:  "eq",
						Value:     "dell",
					},
					{
						Namespace: "sh.hollow.firmware_set.labels",
						Keys:      []string{"model"},
						Operator:  "eq",
						Value:     "r640",
					},
				},
			},
			[]*models.AttributesFirmwareSet{dbtools.FixtureFirmwareSetR640Attribute},
			[]string{"R640"},
			1,
			1,
			false,
			"",
		},
		{
			"list firmware set by attribute params with OR on attribute",
			&fleetdbapi.ComponentFirmwareSetListParams{
				AttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: "sh.hollow.firmware_set.labels",
						Keys:      []string{"model"},
						Operator:  "eq",
						Value:     "r640",
					},
					{
						Namespace:         "sh.hollow.firmware_set.labels",
						Keys:              []string{"model"},
						Operator:          "eq",
						Value:             "x11dph-t",
						AttributeOperator: fleetdbapi.AttributeLogicalOR,
					},
				},
			},
			[]*models.AttributesFirmwareSet{
				dbtools.FixtureFirmwareSetR640Attribute,
				dbtools.FixtureFirmwareSetX11DPHTAttribute,
			},
			[]string{"R640", "X11DPH-T"},
			2,
			1,
			false,
			"",
		},
		{
			"list with incorrect firmware set Name attribute returns no records",
			&fleetdbapi.ComponentFirmwareSetListParams{
				Name: "does-not-exist",
			},
			nil,
			nil,
			0,
			1,
			false,
			"",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			got, resp, err := s.Client.ListServerComponentFirmwareSet(context.TODO(), tt.params)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.NotNil(t, got)

			if tt.expectedFirmwareSetAttribute != nil {
				assert.True(t, assertFirmwareSetAttributeNSEqual(t, tt.expectedFirmwareSetAttribute, got))
				assert.True(t, assertContainsFirmwareSetAttributes(t, tt.expectedFirmwareSetAttribute, got))
			}

			if tt.expectedFirmwareModels != nil {
				assert.True(t, firmwareSetContainsModel(t, tt.expectedFirmwareModels, got))
			}

			assert.Equal(t, tt.expectedPage, resp.Page)
			assert.Equal(t, tt.expectedTotalRecordCount, int(resp.TotalRecordCount))
		})
	}
}

func assertContainsFirmwareSetAttributes(t *testing.T, fwSetModelAttrs []*models.AttributesFirmwareSet, fwSets []fleetdbapi.ComponentFirmwareSet) bool {
	t.Helper()

	expected := len(fwSetModelAttrs)

	var got int

	for _, fwSetModelAttr := range fwSetModelAttrs {
		for _, fwSet := range fwSets {
			if assertAttributesContains(t, fwSet.Attributes, fwSetModelAttr.Data) {
				got++
			}
		}
	}

	return expected == got
}

func assertFirmwareSetAttributeNSEqual(t *testing.T, fwSetModelAttrs []*models.AttributesFirmwareSet, fwSets []fleetdbapi.ComponentFirmwareSet) bool {
	for _, fwSetModelAttr := range fwSetModelAttrs {
		for _, fwSet := range fwSets {
			for _, attr := range fwSet.Attributes {
				if fwSetModelAttr.Namespace != attr.Namespace {
					t.Errorf("attr namespace %s != %s", fwSetModelAttr.Namespace, attr.Namespace)
					return false
				}
			}
		}
	}

	return true
}

func firmwareSetContainsModel(t *testing.T, models []string, set []fleetdbapi.ComponentFirmwareSet) bool {
	t.Helper()

	for _, model := range models {
		for _, f := range set {
			for _, firmware := range f.ComponentFirmware {
				if slices.Contains(firmware.Model, model) {
					return true
				}
			}
		}
	}

	return false
}

func TestIntegrationServerComponentFirmwareSetDelete(t *testing.T) {
	s := serverTest(t)

	var firmwareSetID uuid.UUID

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		var err error

		_, err = s.Client.DeleteServerComponentFirmwareSet(ctx, firmwareSetID)
		if !expectError {
			return nil
		}

		return err
	})

	firmwareSetID, err := uuid.Parse(dbtools.FixtureFirmwareSetR640.ID)
	if err != nil {
		t.Fatal(err)
	}

	var testCases = []struct {
		testName         string
		firmwareSetID    uuid.UUID
		expectedError    bool
		errorMsg         string
		expectedResponse string
	}{
		{
			"component firmware set UUID required",
			uuid.Nil,
			true,
			"",
			"expected a valid firmware set UUID",
		},
		{
			"unknown firmware set UUID returns not found",
			uuid.New(),
			true,
			"",
			"resource not found",
		},
		{
			"firmware set removed",
			firmwareSetID,
			false,
			"",
			"resource deleted",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			resp, err := s.Client.DeleteServerComponentFirmwareSet(context.TODO(), tt.firmwareSetID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Contains(t, err.Error(), tt.expectedResponse)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Contains(t, tt.expectedResponse, resp.Message)
		})
	}
}

func TestIntegrationServerComponentFirmwareSetRemoveFirmware(t *testing.T) {
	s := serverTest(t)

	var firmwareSetID uuid.UUID

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		var err error

		_, err = s.Client.RemoveServerComponentFirmwareSetFirmware(ctx, firmwareSetID, fleetdbapi.ComponentFirmwareSetRequest{})
		if !expectError {
			return nil
		}

		return err
	})

	firmwareSetID, err := uuid.Parse(dbtools.FixtureFirmwareSetR640.ID)
	if err != nil {
		t.Fatal(err)
	}

	// retrieve component firmware fixture data for test
	firmware, _, err := s.Client.GetServerComponentFirmwareSet(context.Background(), firmwareSetID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, firmware)

	r640FirmwareIDs := r640FirmwareFixtureUUIDs(t, firmware.ComponentFirmware)

	var testCases = []struct {
		testName           string
		firmwareSetPayload *fleetdbapi.ComponentFirmwareSetRequest
		expectedError      bool
		errorMsg           string
		expectedResponse   string
	}{
		{
			"component firmware set UUID required",
			&fleetdbapi.ComponentFirmwareSetRequest{},
			true,
			"",
			"expected a valid firmware set UUID",
		},
		{
			"payload must include a non-nil firmware set UUID",
			&fleetdbapi.ComponentFirmwareSetRequest{
				ID:   uuid.Nil,
				Name: "foobar",
			},
			true,
			"",
			"expected a valid firmware set UUID",
		},
		{
			"firmware for removal must be part of firmware set",
			&fleetdbapi.ComponentFirmwareSetRequest{
				ID:                     firmwareSetID,
				ComponentFirmwareUUIDs: []string{uuid.NewString()},
			},
			true,
			"",
			"does not contain firmware",
		},
		{
			"firmware removed from set",
			&fleetdbapi.ComponentFirmwareSetRequest{
				ID:                     firmwareSetID,
				ComponentFirmwareUUIDs: []string{r640FirmwareIDs[0]},
			},
			false,
			"",
			"resource deleted",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			resp, err := s.Client.RemoveServerComponentFirmwareSetFirmware(context.TODO(), tt.firmwareSetPayload.ID, *tt.firmwareSetPayload)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Contains(t, err.Error(), tt.expectedResponse)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Contains(t, tt.expectedResponse, resp.Message)
		})
	}
}
