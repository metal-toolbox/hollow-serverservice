package hollow_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.metalkube.net/hollow/internal/db"
	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestIntegrationHWComponentTypeServiceCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		hct := hollow.HardwareComponentType{Name: "integration-test"}

		res, err := s.Client.HardwareComponentType.Create(ctx, hct)
		if !expectError {
			assert.NotEqual(t, "00000000-0000-0000-0000-000000000000", res.String())
		}

		return err
	})
}

func TestIntegrationHWComponentTypeServiceList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		res, err := s.Client.HardwareComponentType.List(ctx, nil)
		if !expectError {
			assert.Len(t, res, 1)
			assert.Equal(t, db.FixtureHCTFins.ID, res[0].UUID)
			assert.Equal(t, db.FixtureHCTFins.Name, res[0].Name)
		}

		return err
	})
}
