package hollow_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestHardwareServiceListVersionedAttributess(t *testing.T) {
	ctx := context.Background()

	d := time.Now().Add(1 * time.Millisecond)

	timeCtx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	testUUID := uuid.New()

	var testCases = []struct {
		testName     string
		uuid         uuid.UUID
		bios         hollow.VersionedAttributes
		ctx          context.Context
		responseCode int
		expectError  bool
		errorMsg     string
	}{
		{
			"happy path",
			testUUID,
			hollow.VersionedAttributes{EntityUUID: testUUID},
			ctx,
			http.StatusOK,
			false,
			"",
		},
		{
			"server returns an error",
			testUUID,
			hollow.VersionedAttributes{EntityUUID: uuid.New()},
			ctx,
			http.StatusUnauthorized,
			true,
			"server error - response code: 401, message:",
		},
		{
			"fake timeout",
			testUUID,
			hollow.VersionedAttributes{EntityUUID: uuid.New()},
			timeCtx,
			http.StatusOK,
			true,
			"context deadline exceeded",
		},
	}

	for _, tt := range testCases {
		jsonResponse, err := json.Marshal([]hollow.VersionedAttributes{tt.bios})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), tt.responseCode)

		res, err := c.Hardware.GetVersionedAttributes(tt.ctx, tt.uuid)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, res, tt.testName)
		}
	}
}
