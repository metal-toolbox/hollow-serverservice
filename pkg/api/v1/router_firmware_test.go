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

func TestIntegrationFirmwareList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		params := serverservice.ComponentFirmwareVersionListParams{
			Vendor:  "",
			Model:   "",
			Version: "",
		}

		r, resp, err := s.Client.ListFirmware(ctx, &params)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, r, 3)
			assert.EqualValues(t, 3, resp.PageCount)
			assert.EqualValues(t, 1, resp.TotalPages)
			assert.EqualValues(t, 3, resp.TotalRecordCount)
			// We returned everything, so we shouldnt have a next page info
			assert.Nil(t, resp.Links.Next)
			assert.Nil(t, resp.Links.Previous)
		}
		return err
	})

	var testCases = []struct {
		testName      string
		params        *serverservice.ComponentFirmwareVersionListParams
		expectedUUIDs []string
		expectError   bool
		errorMsg      string
	}{
		{
			"search by vendor",
			&serverservice.ComponentFirmwareVersionListParams{
				Vendor: "Dell",
			},
			[]string{dbtools.FixtureDellR640.ID, dbtools.FixtureDellR6515.ID},
			false,
			"",
		},
		{
			"search by model",
			&serverservice.ComponentFirmwareVersionListParams{
				Model: "X11DPH-T",
			},
			[]string{dbtools.FixtureSuperMicro.ID},
			false,
			"",
		},
		{
			"search by version",
			&serverservice.ComponentFirmwareVersionListParams{
				Version: "2.6.6",
			},
			[]string{dbtools.FixtureDellR6515.ID},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, _, err := s.Client.ListFirmware(context.TODO(), tt.params)
			if tt.expectError {
				assert.NoError(t, err)
				return
			}

			var actual []string

			for _, srv := range r {
				actual = append(actual, srv.UUID.String())
			}

			assert.ElementsMatch(t, tt.expectedUUIDs, actual)
		})
	}
}

func TestIntegrationFirmwareGet(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)
		fw, _, err := s.Client.GetFirmware(ctx, uuid.MustParse(dbtools.FixtureDellR640.ID))

		if !expectError {
			require.NoError(t, err)
			assert.Equal(t, fw.UUID, uuid.MustParse(dbtools.FixtureDellR640.ID))
		}

		return err
	})
}

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
			UpstreamURL: "https://vendor.com/firmwares/DSU_21.07.00/",
			S3URL:       "http://example-firmware-bucket.s3.amazonaws.com/firmware/dell/DSU_21.07.00/",
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
		_, err := s.Client.DeleteFirmware(ctx, serverservice.ComponentFirmwareVersion{UUID: uuid.MustParse(dbtools.FixtureDellR640.ID)})

		return err
	})
}

func TestIntegrationFirmwareUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		resp, err := s.Client.UpdateFirmware(ctx, uuid.MustParse(dbtools.FixtureDellR640.ID), serverservice.ComponentFirmwareVersion{Filename: "foobarino"})
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/firmwares/%s", dbtools.FixtureDellR640.ID), resp.Links.Self.Href)
			fw, _, _ := s.Client.GetFirmware(ctx, uuid.MustParse(dbtools.FixtureDellR640.ID))
			assert.Equal(t, "foobarino", fw.Filename)
		}

		return err
	})
}
