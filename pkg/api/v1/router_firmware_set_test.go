package serverservice_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.hollow.sh/serverservice/internal/dbtools"
	"go.hollow.sh/serverservice/internal/models"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

// r640FirmwareFixtureUUIDs returns firmware  uuids based on the firmware hardware model attribute
func r640FirmwareFixtureUUIDs(t *testing.T, firmware []serverservice.ComponentFirmwareVersion) []string {
	t.Helper()

	ids := []string{}

	for idx, f := range firmware {
		if f.Model == "R640" {
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

		var testFirmwareSet serverservice.ComponentFirmwareSetRequest

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

			testFirmwareSet = serverservice.ComponentFirmwareSetRequest{
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
		firmwareSetPayload *serverservice.ComponentFirmwareSetRequest
		expectedError      bool
		expectedResponse   string
		errorMsg           string
	}{
		{
			"Name field required",
			&serverservice.ComponentFirmwareSetRequest{},
			true,
			"400",
			"required attribute not set: Name",
		},
		{
			"component firmware UUIDs required",
			&serverservice.ComponentFirmwareSetRequest{
				Name: "foobar",
			},
			true,
			"400",
			"expected one or more firmware UUIDs, got none",
		},
		{
			"valid UUIDs for the firmware ID expected",
			&serverservice.ComponentFirmwareSetRequest{
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
			&serverservice.ComponentFirmwareSetRequest{
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
			&serverservice.ComponentFirmwareSetRequest{
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
			&serverservice.ComponentFirmwareSetRequest{
				Name:                   "foobar",
				ComponentFirmwareUUIDs: r640FirmwareIDs,
				Attributes: []serverservice.Attributes{
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

		_, err = s.Client.UpdateComponentFirmwareSetRequest(ctx, firmwareSetID, serverservice.ComponentFirmwareSetRequest{})
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
		testName           string
		firmwareSetPayload *serverservice.ComponentFirmwareSetRequest
		expectedError      bool
		expectedResponse   string
		errorMsg           string
	}{

		{
			"component firmware set UUID required",
			&serverservice.ComponentFirmwareSetRequest{
				Name: "foobar",
			},
			true,
			"400",
			"expected a valid firmware set ID, got none",
		},
		{
			"valid UUIDs for the firmware ID expected",
			&serverservice.ComponentFirmwareSetRequest{
				ID:   firmwareSetID,
				Name: "foobar",
				ComponentFirmwareUUIDs: []string{
					"d825bbeb-20fb-452e-9fe4-invalid",
				},
			},
			true,
			"400",
			"invalid firmware UUID",
		},
		{
			"duplicate firmware UUIDs are not accepted",
			&serverservice.ComponentFirmwareSetRequest{
				ID:   firmwareSetID,
				Name: "foobar",
				ComponentFirmwareUUIDs: []string{
					r640FirmwareIDs[0],
					r640FirmwareIDs[0],
				},
			},
			true,
			"400",
			"exists in firmware set",
		},
		{
			"non-existing firmware UUID are not accepted",
			&serverservice.ComponentFirmwareSetRequest{
				ID:   firmwareSetID,
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
			"update an existing firmware set - update name, referenced firmware",
			&serverservice.ComponentFirmwareSetRequest{
				Name:                   "foobar-updated",
				ID:                     firmwareSetID,
				ComponentFirmwareUUIDs: []string{dbtools.FixtureDellR640CPLD.ID},
			},
			false,
			"200",
			"",
		},
		{
			"update an existing firmware set - update metadata",
			&serverservice.ComponentFirmwareSetRequest{
				Name: "foobar-updated",
				ID:   firmwareSetID,
				Attributes: []serverservice.Attributes{
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
			assert.Equal(t, 3, len(got.ComponentFirmware))

			// assert firmware set attributes
			if len(tt.firmwareSetPayload.Attributes) > 0 {
				assert.Equal(t, 1, len(got.Attributes))
				assert.Equal(t, tt.firmwareSetPayload.Attributes[0].Namespace, got.Attributes[0].Namespace)
				assertAttributesEqual(t, tt.firmwareSetPayload.Attributes[0].Data, got.Attributes[0].Data)
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
			assertAttributesEqual(t, []byte(dbtools.FixtureFirmwareSetR640Attribute.Data), got.Attributes[0].Data)

			// assert component firmware
			assert.Equal(t, 2, len(got.ComponentFirmware))
			assert.Equal(t, "r640", got.Name)
		})
	}
}

func assertAttributesEqual(t *testing.T, a, b []byte) {
	t.Helper()

	// unmarshal fixture attribute data
	aData := map[string]string{}
	if err := json.Unmarshal(a, &aData); err != nil {
		t.Fatal(err)
	}

	// unmarshal got attribute data
	bData := map[string]string{}
	if err := json.Unmarshal(b, &bData); err != nil {
		t.Fatal(err)
	}

	// assert fixture data
	assert.Equal(t, aData, bData)
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
		params                       *serverservice.ComponentFirmwareSetListParams
		expectedFirmwareSetAttribute *models.Attribute
		expectedFirmwareModel        string
		expectedFirmwareCount        int
		expectedTotalRecordCount     int
		expectedPage                 int
		expectedError                bool
		errorMsg                     string
	}{

		{
			"list firmware set by name - r640",
			&serverservice.ComponentFirmwareSetListParams{Name: "r640"},
			dbtools.FixtureFirmwareSetR640Attribute,
			"R640",
			2,
			1,
			1,
			false,
			"",
		},
		{
			"list firmware set by name - r6515",
			&serverservice.ComponentFirmwareSetListParams{Name: "r6515"},
			dbtools.FixtureFirmwareSetR6515Attribute,
			"R6515",
			2,
			1,
			1,
			false,
			"",
		},
		{
			"list with pagination Limit, Offset",
			&serverservice.ComponentFirmwareSetListParams{
				Pagination: &serverservice.PaginationParams{
					Limit: 1,
					Page:  2,
				},
			},
			nil,
			"",
			2,
			2,
			2,
			false,
			"",
		},
		{
			"list with incorrect firmware set Name attribute returns no records",
			&serverservice.ComponentFirmwareSetListParams{
				Name: "does-not-exist",
			},
			nil,
			"",
			0,
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

			if tt.expectedFirmwareCount > 0 {
				assert.Equal(t, tt.expectedFirmwareCount, len(got[0].ComponentFirmware))
			}

			if tt.expectedFirmwareSetAttribute != nil {
				assert.Equal(t, 1, len(got[0].Attributes))
				assert.Equal(t, tt.expectedFirmwareSetAttribute.Namespace, got[0].Attributes[0].Namespace)
				assertAttributesEqual(t, []byte(tt.expectedFirmwareSetAttribute.Data), got[0].Attributes[0].Data)
			}

			if tt.expectedFirmwareModel != "" {
				assert.Equal(t, tt.expectedFirmwareModel, got[0].ComponentFirmware[0].Model)
				assert.Equal(t, tt.expectedFirmwareModel, got[0].ComponentFirmware[1].Model)
			}

			assert.Equal(t, tt.expectedPage, resp.Page)
			assert.Equal(t, tt.expectedTotalRecordCount, int(resp.TotalRecordCount))
		})
	}
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

		_, err = s.Client.RemoveServerComponentFirmwareSetFirmware(ctx, firmwareSetID, serverservice.ComponentFirmwareSetRequest{})
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
		firmwareSetPayload *serverservice.ComponentFirmwareSetRequest
		expectedError      bool
		errorMsg           string
		expectedResponse   string
	}{
		{
			"component firmware set UUID required",
			&serverservice.ComponentFirmwareSetRequest{},
			true,
			"",
			"expected a valid firmware set UUID",
		},
		{
			"payload must include a non-nil firmware set UUID",
			&serverservice.ComponentFirmwareSetRequest{
				ID:   uuid.Nil,
				Name: "foobar",
			},
			true,
			"",
			"expected a valid firmware set UUID",
		},
		{
			"firmware for removal must be part of firmware set",
			&serverservice.ComponentFirmwareSetRequest{
				ID:                     firmwareSetID,
				ComponentFirmwareUUIDs: []string{uuid.NewString()},
			},
			true,
			"",
			"does not contain firmware",
		},
		{
			"firmware removed from set",
			&serverservice.ComponentFirmwareSetRequest{
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
