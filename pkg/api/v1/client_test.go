package hollow_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestNewClient(t *testing.T) {
	var testCases = []struct {
		testName    string
		authToken   string
		url         string
		expectError bool
		errorMsg    string
	}{
		{
			"no authToken",
			"",
			"https://hollow.metalkube.net",
			true,
			"failed to initialize: no auth token provided",
		},
		{
			"no uri",
			"SuperSecret",
			"",
			true,
			"failed to initialize: no hollow api url provided",
		},
		{
			"happy path",
			"SuperSecret",
			"https://hollow.metalkube.net",
			false,
			"",
		},
	}

	for _, tt := range testCases {
		c, err := hollow.NewClient(tt.authToken, tt.url, nil)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, c, tt.testName)
			assert.NotNil(t, c.Hardware, tt.testName)
		}
	}
}
