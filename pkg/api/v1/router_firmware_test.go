package serverservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.hollow.sh/serverservice/internal/dbtools"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

func TestIntegrationFirmwareCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		testFirmware := serverservice.ComponentFirmwareVersion{
			UUID:        uuid.New(),
			Vendor:      "Dell",
			Model:       "R615",
			Filename:    "foobar",
			Version:     "21.07.00",
			Component:   "system",
			Utility:     "dsu",
			Sha:         "foobar",
			UpstreamURL: "https://linux.dell.com/repo/hardware/DSU_21.07.00/",
		}

		id, resp, err := s.Client.CreateFirmware(ctx, testFirmware)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, id)
			assert.Equal(t, testFirmware.UUID.String(), id.String())
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/firmwares/%s", id), resp.Links.Self.Href)
		}

		return err
	})
}

func TestIntegrationFirmwareDelete(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)
		_, err := s.Client.DeleteFirmware(ctx, serverservice.ComponentFirmwareVersion{UUID: uuid.MustParse(dbtools.FixtureDell210700.ID)})

		return err
	})
}

func TestIntegrationFirmwareUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		resp, err := s.Client.UpdateFirmware(ctx, uuid.MustParse(dbtools.FixtureDell210700.ID), serverservice.ComponentFirmwareVersion{Filename: "foobarino"})
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/firmwares/%s", dbtools.FixtureDell210700.ID), resp.Links.Self.Href)
			fw, _, _ := s.Client.GetFirmware(ctx, uuid.MustParse(dbtools.FixtureDell210700.ID))
			assert.Equal(t, "foobarino", fw.Filename)
		}

		return err
	})
}
