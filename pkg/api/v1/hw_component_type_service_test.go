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

func TestHWComponentTypeServiceCreate(t *testing.T) {
	ctx := context.Background()

	d := time.Now().Add(1 * time.Millisecond)

	timeCtx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	var testCases = []struct {
		testName     string
		hct          hollow.HardwareComponentType
		ctx          context.Context
		responseCode int
		expectError  bool
		errorMsg     string
	}{
		{
			"happy path",
			hollow.HardwareComponentType{Name: "Test"},
			ctx,
			http.StatusOK,
			false,
			"",
		},
		{
			"server returns an error",
			hollow.HardwareComponentType{Name: "Test"},
			ctx,
			http.StatusUnauthorized,
			true,
			"server error - response code: 401, message: something something",
		},
		{
			"fake timeout",
			hollow.HardwareComponentType{Name: "Test"},
			timeCtx,
			http.StatusOK,
			true,
			"context deadline exceeded",
		},
	}

	for _, tt := range testCases {
		c := mockClient(`{"message": "something something"}`, tt.responseCode)

		err := c.HardwareComponentType.Create(tt.ctx, tt.hct)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, c, tt.testName)
		}
	}
}

func TestHWComponentTypeServiceList(t *testing.T) {
	ctx := context.Background()

	d := time.Now().Add(1 * time.Millisecond)

	timeCtx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	testUUID := uuid.New()

	var testCases = []struct {
		testName     string
		filter       *hollow.HardwareComponentTypeListParams
		hct          hollow.HardwareComponentType
		ctx          context.Context
		responseCode int
		expectError  bool
		errorMsg     string
	}{
		{
			"happy path - no filter",
			&hollow.HardwareComponentTypeListParams{},
			hollow.HardwareComponentType{UUID: testUUID, Name: "Test-CPU"},
			ctx,
			http.StatusOK,
			false,
			"",
		},
		{
			"happy path - filter",
			&hollow.HardwareComponentTypeListParams{Name: "Test-CPU"},
			hollow.HardwareComponentType{UUID: testUUID, Name: "Test-CPU"},
			ctx,
			http.StatusOK,
			false,
			"",
		},
		{
			"server returns an error",
			nil,
			hollow.HardwareComponentType{UUID: testUUID, Name: "Test-CPU"},
			ctx,
			http.StatusUnauthorized,
			true,
			"server error - response code: 401, message:",
		},
		{
			"fake timeout",
			nil,
			hollow.HardwareComponentType{UUID: testUUID, Name: "Test-CPU"},
			timeCtx,
			http.StatusOK,
			true,
			"context deadline exceeded",
		},
	}

	for _, tt := range testCases {
		jsonResponse, err := json.Marshal([]hollow.HardwareComponentType{tt.hct})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), tt.responseCode)

		res, err := c.HardwareComponentType.List(tt.ctx, tt.filter)

		if tt.expectError {
			assert.Error(t, err, tt.testName)
			assert.Contains(t, err.Error(), tt.errorMsg)
		} else {
			assert.NoError(t, err, tt.testName)
			assert.NotNil(t, res, tt.testName)
		}
	}
}
