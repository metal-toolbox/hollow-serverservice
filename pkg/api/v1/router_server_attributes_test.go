package hollow_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"go.metalkube.net/hollow/internal/db"
)

func TestIntegrationServerAttributesUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		err := s.Client.Server.UpdateAttributes(ctx, db.FixtureServerDory.ID, db.FixtureNamespaceMetadata, json.RawMessage([]byte(`{"setting":"enabled"}`)))
		if !expectError {
			// assert.Nil(t, resp.Links.Self)
			require.NoError(t, err)
		}

		return err
	})
}
