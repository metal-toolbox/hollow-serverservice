package hollow_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.metalkube.net/hollow/internal/db"
	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestIntegrationServerComponentTypeServiceCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		hct := hollow.ServerComponentType{Name: "integration-test"}

		r, resp, err := s.Client.ServerComponentType.Create(ctx, hct)
		if !expectError {
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil.String(), r.String())
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("/api/v1/server-component-types/%s", r.String()), resp.Links.Self.Href)
		}

		return err
	})
}

func TestIntegrationServerComponentTypeServiceList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		r, resp, err := s.Client.ServerComponentType.List(ctx, nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, r, 1)
			assert.Equal(t, db.FixtureSCTFins.ID, r[0].UUID)
			assert.Equal(t, db.FixtureSCTFins.Name, r[0].Name)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Links.Self)
		}

		return err
	})
}
