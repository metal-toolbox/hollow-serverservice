package serverservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.hollow.sh/serverservice/internal/dbtools"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

func TestIntegrationCreateServerConditionType(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		hct := &serverservice.ServerConditionType{Slug: "integration-test"}

		resp, err := s.Client.CreateServerConditionType(ctx, hct)
		if !expectError {
			require.NoError(t, err)
			assert.Equal(t, "integration-test", resp.Slug)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/server-condition-types/%s", resp.Slug), resp.Links.Self.Href)
		}

		return err
	})

	t.Run("creating a duplicate server condition type fails", func(t *testing.T) {
		_, err := s.Client.CreateServerConditionType(context.TODO(), &serverservice.ServerConditionType{Slug: dbtools.FixtureSwimConditionType.Slug})
		assert.Error(t, err)
		require.Contains(t, err.Error(), "duplicate key")
	})
}

func TestIntegrationListServerConditionTypes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		r, resp, err := s.Client.ListServerConditionTypes(ctx, nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, r, 1)
			assert.Equal(t, dbtools.FixtureSwimConditionType.Slug, r[0].Slug)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Links.Self)
		}

		return err
	})
}
