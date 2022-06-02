package serverservice_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hollow "go.hollow.sh/serverservice/pkg/api/v1"
)

func TestServerComponentTypeServiceCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		hct := hollow.ServerComponentType{Name: "unit-test"}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created", "slug":"slug-1"}`))

		c := mockClient(string(jsonResponse), respCode)
		resp, err := c.CreateServerComponentType(ctx, hct)
		if !expectError {
			assert.Equal(t, "slug-1", resp.Slug)
		}

		return err
	})
}

func TestServerComponentTypeServiceList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		hct := hollow.ServerComponentTypeSlice{{Slug: "slug-1", Name: "unit-test-1"}, {Slug: "slug-2", Name: "unit-test-2"}}
		jsonResponse, err := json.Marshal(hollow.ServerResponse{Records: hct})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.ListServerComponentTypes(ctx, nil)
		if !expectError {
			assert.ElementsMatch(t, hct, res)
		}

		return err
	})
}
