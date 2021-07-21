package hollow_test

import (
	"context"
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

		res, err := s.Client.ServerComponentType.Create(ctx, hct)
		if !expectError {
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil.String(), res.String())
		}

		return err
	})
}

func TestIntegrationServerComponentTypeServiceList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		res, err := s.Client.ServerComponentType.List(ctx, nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, res, 1)
			assert.Equal(t, db.FixtureSCTFins.ID, res[0].UUID)
			assert.Equal(t, db.FixtureSCTFins.Name, res[0].Name)
		}

		return err
	})
}
