package hollow_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.metalkube.net/hollow/internal/gormdb"
)

func TestIntegrationServerListComponents(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		attrs, _, err := s.Client.Server.ListComponents(ctx, gormdb.FixtureServerNemo.ID, nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, attrs, 2)
		}

		return err
	})

	var testCases = []struct {
		testName string
		srvUUID  uuid.UUID
		errorMsg string
	}{
		{
			"returns not found on missing server uuid",
			uuid.New(),
			"response code: 404",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			_, _, err := s.Client.Server.ListComponents(context.TODO(), tt.srvUUID, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}
