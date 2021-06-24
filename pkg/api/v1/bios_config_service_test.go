package hollow_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestBiosConfigServiceCreateBIOSConfig(t *testing.T) {
	ctx := context.Background()

	exampleBiosResults := `{
    "open": {
      "boot_mode": "Bios"
    }
  }`

	jsonBios, err := json.Marshal(exampleBiosResults)
	if err != nil {
		fmt.Println("failed to convert example bios to json")
		log.Fatal(err)
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
	}

	for _, tt := range testCases {
		c := mockClient("", tt.responseCode)

		err := c.BIOSConfig.CreateBIOSConfig(ctx, tt.bios)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, c, tt.testName)
		}
	}
}
