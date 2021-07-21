package hollow_test

import (
	"context"
	"net/http"
	"testing"
	"time"

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
			assert.NotNil(t, c.Server, tt.testName)
			assert.NotNil(t, c.ServerComponentType, tt.testName)
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
			validToken([]string{"read", "write"}),
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
			http.StatusUnauthorized,
			true,
			"server error - response code: 401, message:",
		},
		{
			"fake timeout",
			timeCtx,
			validToken([]string{"read", "write"}),
			http.StatusOK,
			true,
			"context deadline exceeded",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
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
