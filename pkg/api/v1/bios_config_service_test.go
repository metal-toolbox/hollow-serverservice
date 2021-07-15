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
			"server error - response code: 401, message: something something",
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
		c := mockClient(`{"message": "something something", "uuid":"00000000-0000-0000-0000-000000001234"}`, tt.responseCode)

		r, err := c.BIOSConfig.Create(tt.ctx, tt.bios)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, c, tt.testName)
			assert.NotNil(t, r)
			assert.Equal(t, "00000000-0000-0000-0000-000000001234", r.String())
		}
	}
}
