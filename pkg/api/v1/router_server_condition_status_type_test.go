package serverservice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"go.hollow.sh/serverservice/internal/dbtools"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

func TestIntegrationCreateServerConditionStatusType(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		hct := &serverservice.ServerConditionStatusType{Slug: "integration-test"}

		resp, err := s.Client.CreateServerConditionStatusType(ctx, hct)
		if !expectError {
			require.NoError(t, err)
			assert.Equal(t, "integration-test", resp.Slug)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, "http://test.hollow.com/api/v1/server-condition-status-types/"+resp.Slug, resp.Links.Self.Href)
		}

		return err
	})

	s.Client.SetToken(validToken(adminScopes))

	t.Run("creating a duplicate server condition status type fails", func(t *testing.T) {
		_, err := s.Client.CreateServerConditionStatusType(context.TODO(), &serverservice.ServerConditionStatusType{Slug: dbtools.FixtureSwimConditionStatusTypeActive.Slug})
		assert.Error(t, err)
		require.Contains(t, err.Error(), "duplicate key")
	})
}

func TestIntegrationListServerConditionStatusTypes(t *testing.T) {
	s := serverTest(t)

	fixtureConditions := []string{dbtools.FixtureSwimConditionStatusTypeActive.Slug, dbtools.FixtureSwimConditionStatusTypeSucceeded.Slug}

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		r, resp, err := s.Client.ListServerConditionStatusTypes(ctx, nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, r, 2)
			assert.True(t, slices.Contains(fixtureConditions, r[0].Slug))
			assert.True(t, slices.Contains(fixtureConditions, r[1].Slug))
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, "http://test.hollow.com/api/v1/server-condition-status-types", resp.Links.Self.Href)
		}

		return err
	})
}
