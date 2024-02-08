package fleetdbapi_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

func TestServerServiceCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		srv := fleetdbapi.Server{UUID: uuid.New(), FacilityCode: "Test1"}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created", "slug":"00000000-0000-0000-0000-000000001234"}`))

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.Create(ctx, srv)
		if !expectError {
			assert.Equal(t, "00000000-0000-0000-0000-000000001234", res.String())
		}

		return err
	})
}

func TestServerServiceDelete(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse := json.RawMessage([]byte(`{"message": "resource deleted"}`))
		c := mockClient(string(jsonResponse), respCode)
		_, err := c.Delete(ctx, fleetdbapi.Server{UUID: uuid.New()})

		return err
	})
}
func TestServerServiceGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		srv := fleetdbapi.Server{UUID: uuid.New(), FacilityCode: "Test1"}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Record: srv})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.Get(ctx, srv.UUID)
		if !expectError {
			assert.Equal(t, srv.UUID, res.UUID)
			assert.Equal(t, srv.FacilityCode, res.FacilityCode)
		}

		return err
	})
}

func TestServerServiceList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		srv := []fleetdbapi.Server{{UUID: uuid.New(), FacilityCode: "Test1"}}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: srv})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.List(ctx, nil)
		if !expectError {
			assert.ElementsMatch(t, srv, res)
		}

		return err
	})
}

func TestServerServiceUpdate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource updated"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		_, err = c.Update(ctx, uuid.UUID{}, fleetdbapi.Server{Name: "new-name"})

		return err
	})
}

func TestServerServiceCreateAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		attr := fleetdbapi.Attributes{Namespace: "unit-test", Data: json.RawMessage([]byte(`{"test":"unit"}`))}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created"}`))

		c := mockClient(string(jsonResponse), respCode)
		_, err := c.CreateAttributes(ctx, uuid.New(), attr)

		return err
	})
}
func TestServerServiceGetAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		attr := &fleetdbapi.Attributes{Namespace: "unit-test", Data: json.RawMessage([]byte(`{"test":"unit"}`))}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Record: attr})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.GetAttributes(ctx, uuid.UUID{}, "unit-test")
		if !expectError {
			assert.Equal(t, attr, res)
		}

		return err
	})
}

func TestServerServiceDeleteAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource deleted"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		_, err = c.DeleteAttributes(ctx, uuid.UUID{}, "unit-test")

		return err
	})
}

func TestServerServiceListAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		attrs := []fleetdbapi.Attributes{{Namespace: "unit-test", Data: json.RawMessage([]byte(`{"test":"unit"}`))}}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: attrs})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.ListAttributes(ctx, uuid.UUID{}, nil)
		if !expectError {
			assert.ElementsMatch(t, attrs, res)
		}

		return err
	})
}

func TestServerServiceUpdateAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource updated"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		_, err = c.UpdateAttributes(ctx, uuid.UUID{}, "unit-test", json.RawMessage([]byte(`{"test":"unit"}`)))

		return err
	})
}

func TestServerServiceComponentsGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		sc := []fleetdbapi.ServerComponent{{Name: "unit-test", Serial: "1234"}}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: sc})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.GetComponents(ctx, uuid.UUID{}, nil)
		if !expectError {
			assert.ElementsMatch(t, sc, res)
		}

		return err
	})
}

func TestServerServiceComponentsList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		sc := []fleetdbapi.ServerComponent{{Name: "unit-test", Serial: "1234"}}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: sc})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.ListComponents(ctx, &fleetdbapi.ServerComponentListParams{Name: "unit-test", Serial: "1234"})
		if !expectError {
			assert.ElementsMatch(t, sc, res)
		}

		return err
	})
}

func TestServerServiceComponentsCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource created"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.CreateComponents(ctx, uuid.New(), fleetdbapi.ServerComponentSlice{{Name: "unit-test"}})
		if !expectError {
			assert.Contains(t, res.Message, "resource created")
		}

		return err
	})
}

func TestServerServiceComponentsUpdate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource updated"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.UpdateComponents(ctx, uuid.New(), fleetdbapi.ServerComponentSlice{{Name: "unit-test"}})
		if !expectError {
			assert.Contains(t, res.Message, "resource updated")
		}

		return err
	})
}

func TestServerServiceVersionedAttributeCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		va := fleetdbapi.VersionedAttributes{Namespace: "unit-test", Data: json.RawMessage([]byte(`{"test":"unit"}`))}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created", "slug":"the-namespace"}`))

		c := mockClient(string(jsonResponse), respCode)
		resp, err := c.CreateVersionedAttributes(ctx, uuid.New(), va)
		if !expectError {
			assert.Equal(t, "the-namespace", resp.Slug)
		}

		return err
	})
}

func TestServerServiceGetVersionedAttributess(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		va := []fleetdbapi.VersionedAttributes{{Namespace: "test", Data: json.RawMessage([]byte(`{}`))}}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: va})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.GetVersionedAttributes(ctx, uuid.New(), "namespace")
		if !expectError {
			assert.ElementsMatch(t, va, res)
		}

		return err
	})
}

func TestServerServiceListVersionedAttributess(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		va := []fleetdbapi.VersionedAttributes{{Namespace: "test", Data: json.RawMessage([]byte(`{}`))}}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: va})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.ListVersionedAttributes(ctx, uuid.New())
		if !expectError {
			assert.ElementsMatch(t, va, res)
		}

		return err
	})
}

func TestServerServiceCreateServerComponentFirmware(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		firmware := fleetdbapi.ComponentFirmwareVersion{
			UUID:    uuid.New(),
			Vendor:  "Dell",
			Model:   []string{"R615"},
			Version: "21.07.00",
		}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created", "slug":"00000000-0000-0000-0000-000000001234"}`))

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.CreateServerComponentFirmware(ctx, firmware)
		if !expectError {
			assert.Equal(t, "00000000-0000-0000-0000-000000001234", res.String())
		}

		return err
	})
}

func TestServerServiceServerComponentFirmwareDelete(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse := json.RawMessage([]byte(`{"message": "resource deleted"}`))
		c := mockClient(string(jsonResponse), respCode)
		_, err := c.DeleteServerComponentFirmware(ctx, fleetdbapi.ComponentFirmwareVersion{UUID: uuid.New()})

		return err
	})
}
func TestServerServiceServerComponentFirmwareGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		firmware := fleetdbapi.ComponentFirmwareVersion{
			UUID:    uuid.New(),
			Vendor:  "Dell",
			Model:   []string{"R615"},
			Version: "21.07.00",
		}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Record: firmware})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.GetServerComponentFirmware(ctx, firmware.UUID)
		if !expectError {
			assert.Equal(t, firmware.UUID, res.UUID)
			assert.Equal(t, firmware.Vendor, res.Vendor)
			assert.Equal(t, firmware.Model, res.Model)
			assert.Equal(t, firmware.Version, res.Version)
		}

		return err
	})
}

func TestServerServiceServerComponentFirmwareList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		firmware := []fleetdbapi.ComponentFirmwareVersion{{
			UUID:    uuid.New(),
			Vendor:  "Dell",
			Model:   []string{"R615"},
			Version: "21.07.00",
		}}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: firmware})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.ListServerComponentFirmware(ctx, nil)
		if !expectError {
			assert.ElementsMatch(t, firmware, res)
		}

		return err
	})
}

func TestServerServiceServerComponentFirmwareUpdate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource updated"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		_, err = c.UpdateServerComponentFirmware(ctx, uuid.UUID{}, fleetdbapi.ComponentFirmwareVersion{UUID: uuid.New()})

		return err
	})
}

func TestBillOfMaterialsBatchUpload(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		bom := []fleetdbapi.Bom{{SerialNum: "fakeSerialNum1", AocMacAddress: "fakeAocMacAddress1", BmcMacAddress: "fakeBmcMacAddress1"}}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Record: bom})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.BillOfMaterialsBatchUpload(ctx, bom)
		if !expectError {
			assert.Equal(t, []interface{}([]interface{}{
				map[string]interface{}{
					"aoc_mac_address": "fakeAocMacAddress1",
					"bmc_mac_address": "fakeBmcMacAddress1",
					"metro":           "",
					"num_def_pwd":     "",
					"num_defi_pmi":    "",
					"serial_num":      "fakeSerialNum1"}}), res.Record)
		}

		return err
	})
}

func TestGetBomInfoByAOCMacAddr(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		bom := fleetdbapi.Bom{SerialNum: "fakeSerialNum1", AocMacAddress: "fakeAocMacAddress1", BmcMacAddress: "fakeBmcMacAddress"}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Record: bom})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		respBom, _, err := c.GetBomInfoByAOCMacAddr(ctx, "fakeAocMacAddress1")
		if !expectError {
			assert.Equal(t, &bom, respBom)
		}

		return err
	})
}

func TestGetBomInfoByBMCMacAddr(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		bom := fleetdbapi.Bom{SerialNum: "fakeSerialNum1", AocMacAddress: "fakeAocMacAddress1", BmcMacAddress: "fakeBmcMacAddress1"}
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Record: bom})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		respBom, _, err := c.GetBomInfoByBMCMacAddr(ctx, "fakeBmcMacAddress1")
		if !expectError {
			assert.Equal(t, &bom, respBom)
		}

		return err
	})
}
