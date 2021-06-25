package hollow_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestHardwareServiceListBIOSConfigs(t *testing.T) {
	ctx := context.Background()

	d := time.Now().Add(1 * time.Millisecond)

	timeCtx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	exampleBiosResults := `{
			"open": {
				"boot_mode": "Bios"
			}
		}`

	jsonBios, err := json.Marshal(exampleBiosResults)
	if err != nil {
		fmt.Println("failed to convert example bios to json")
	}

	testUUID := uuid.New()

	var testCases = []struct {
		testName     string
		uuid         uuid.UUID
		bios         hollow.BIOSConfig
		ctx          context.Context
		responseCode int
		expectError  bool
		errorMsg     string
	}{
		{
			"happy path",
			testUUID,
			hollow.BIOSConfig{HardwareUUID: testUUID, ConfigValues: jsonBios},
			ctx,
			http.StatusOK,
			false,
			"",
		},
		{
			"server returns an error",
			testUUID,
			hollow.BIOSConfig{HardwareUUID: uuid.New(), ConfigValues: jsonBios},
			ctx,
			http.StatusUnauthorized,
			true,
			"server error: status_code: 401, message:",
		},
		{
			"fake timeout",
			testUUID,
			hollow.BIOSConfig{HardwareUUID: uuid.New(), ConfigValues: jsonBios},
			timeCtx,
			http.StatusOK,
			true,
			"context deadline exceeded",
		},
	}

	for _, tt := range testCases {
		jsonResponse, err := json.Marshal([]hollow.BIOSConfig{tt.bios})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), tt.responseCode)

		res, err := c.Hardware.ListBIOSConfigs(tt.ctx, tt.uuid)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, res, tt.testName)
		}
	}
}
