package hollow_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestHardwareServiceCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		hw := hollow.Hardware{UUID: uuid.New(), FacilityCode: "Test1"}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created", "uuid":"00000000-0000-0000-0000-000000001234"}`))

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.Hardware.Create(ctx, hw)
		if !expectError {
			assert.Equal(t, "00000000-0000-0000-0000-000000001234", res.String())
		}

		return err
	})
}

func TestHardwareServiceDelete(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse := json.RawMessage([]byte(`{"message": "resource deleted"}`))
		c := mockClient(string(jsonResponse), respCode)

		return c.Hardware.Delete(ctx, hollow.Hardware{UUID: uuid.New()})
	})
}
func TestHardwareServiceGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		hw := hollow.Hardware{UUID: uuid.New(), FacilityCode: "Test1"}
		jsonResponse, err := json.Marshal(hw)
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.Hardware.Get(ctx, hw.UUID)
		if !expectError {
			assert.Equal(t, hw.UUID, res.UUID)
			assert.Equal(t, hw.FacilityCode, res.FacilityCode)
		}

		return err
	})
}

func TestHardwareServiceList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		hw := []hollow.Hardware{{UUID: uuid.New(), FacilityCode: "Test1"}}
		jsonResponse, err := json.Marshal(hw)
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.Hardware.List(ctx, nil)
		if !expectError {
			assert.ElementsMatch(t, hw, res)
		}

		return err
	})
}

func TestHardwareServiceVersionedAttributeCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		va := hollow.VersionedAttributes{Namespace: "unit-test", Values: json.RawMessage([]byte(`{"test":"unit"}`))}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created", "uuid":"00000000-0000-0000-0000-000000001234"}`))

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.Hardware.CreateVersionedAttributes(ctx, uuid.New(), va)
		if !expectError {
			assert.Equal(t, "00000000-0000-0000-0000-000000001234", res.String())
		}

		return err
	})
}

func TestHardwareServiceListVersionedAttributess(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		va := []hollow.VersionedAttributes{{Namespace: "test", Values: json.RawMessage([]byte(`{}`))}}
		jsonResponse, err := json.Marshal(va)
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.Hardware.GetVersionedAttributes(ctx, uuid.New())
		if !expectError {
			assert.ElementsMatch(t, va, res)
		}

		return err
	})
}
