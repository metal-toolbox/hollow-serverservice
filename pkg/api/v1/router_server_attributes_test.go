package hollow_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.metalkube.net/hollow/internal/db"
	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestIntegrationServerAttributesCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		attrs := hollow.Attributes{
			Namespace: "integration.tests",
			Data:      json.RawMessage([]byte(`{"setting":"enabled"}`)),
		}

		resp, err := s.Client.Server.CreateAttributes(ctx, db.FixtureServerNemo.ID, attrs)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/servers/%s/attributes/integration.tests", db.FixtureServerNemo.ID), resp.Links.Self.Href)
			assert.Equal(t, "integration.tests", resp.Slug)
		}

		return err
	})
}

func TestIntegrationServerAttributesUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		_, err := s.Client.Server.UpdateAttributes(ctx, db.FixtureServerDory.ID, db.FixtureNamespaceMetadata, json.RawMessage([]byte(`{"setting":"enabled"}`)))
		if !expectError {
			// assert.Nil(t, resp.Links.Self)
			require.NoError(t, err)
		}

		return err
	})
}
