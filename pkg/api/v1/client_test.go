package serverservice_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	adminScopes = []string{"read", "write", "read:server:secrets", "write:server:secrets"}
)

func TestNewClientWithToken(t *testing.T) {
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
			"https://dcim.hollow.sh",
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
			"https://dcim.hollow.sh",
			false,
			"",
		},
	}

	for _, tt := range testCases {
		c, err := serverservice.NewClientWithToken(tt.authToken, tt.url, nil)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, c, tt.testName)
		}
	}
}

func mockClientTests(t *testing.T, f func(ctx context.Context, respCode int, expectError bool) error) {
	ctx := context.Background()
	timeCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)

	defer cancel()

	var testCases = []struct {
		testName     string
		ctx          context.Context
		responseCode int
		expectError  bool
		errorMsg     string
	}{
		{
			"happy path",
			ctx,
			http.StatusOK,
			false,
			"",
		},
		{
			"server unauthorized",
			ctx,
			http.StatusUnauthorized,
			true,
			"server error - response code: 401, message:",
		},
		{
			"fake timeout",
			timeCtx,
			http.StatusOK,
			true,
			"context deadline exceeded",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			err := f(tt.ctx, tt.responseCode, tt.expectError)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func realClientTests(t *testing.T, f func(ctx context.Context, token string, respCode int, expectError bool) error) {
	ctx := context.Background()
	timeCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)

	defer cancel()

	var testCases = []struct {
		testName     string
		ctx          context.Context
		authToken    string
		responseCode int
		expectError  bool
		errorMsg     string
	}{
		{
			"happy path",
			ctx,
			validToken(adminScopes),
			http.StatusOK,
			false,
			"",
		},
		{
			"invalid auth token",
			ctx,
			"invalidToken",
			http.StatusUnauthorized,
			true,
			"server error - response code: 401, message:",
		},
		// this acts as a safeguard test to ensure that all methods require at least one scope
		{
			"auth token with no scopes",
			ctx,
			validToken([]string{}),
			http.StatusForbidden,
			true,
			"server error - response code: 403, message:",
		},
		{
			"fake timeout",
			timeCtx,
			validToken(adminScopes),
			http.StatusOK,
			true,
			"context deadline exceeded",
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("real client - %s", tt.testName), func(t *testing.T) {
			err := f(tt.ctx, tt.authToken, tt.responseCode, tt.expectError)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
