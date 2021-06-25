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

	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestBiosConfigServiceCreateBIOSConfig(t *testing.T) {
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

	var testCases = []struct {
		testName     string
		bios         hollow.BIOSConfig
		ctx          context.Context
		responseCode int
		expectError  bool
		errorMsg     string
	}{
		{
			"happy path",
			hollow.BIOSConfig{HardwareUUID: uuid.New(), ConfigValues: jsonBios},
			ctx,
			http.StatusOK,
			false,
			"",
		},
		{
			"server returns an error",
			hollow.BIOSConfig{HardwareUUID: uuid.New(), ConfigValues: jsonBios},
			ctx,
			http.StatusUnauthorized,
			true,
			"server error: status_code: 401, message: response body",
		},
		{
			"fake timeout",
			hollow.BIOSConfig{HardwareUUID: uuid.New(), ConfigValues: jsonBios},
			timeCtx,
			http.StatusOK,
			true,
			"context deadline exceeded",
		},
	}

	for _, tt := range testCases {
		c := mockClient("response body", tt.responseCode)

		err := c.BIOSConfig.CreateBIOSConfig(tt.ctx, tt.bios)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, c, tt.testName)
		}
	}
}
