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

	scopes := []string{"read:server-component-firmwares", "write:server-component-firmwares"}
	scopedRealClientTests(t, scopes, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		params := serverservice.ComponentFirmwareVersionListParams{
			Vendor:  "",
			Model:   nil,
			Version: "",
		}

		r, resp, err := s.Client.ListServerComponentFirmware(ctx, &params)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, r, 6)
			assert.EqualValues(t, 6, resp.PageCount)
			assert.EqualValues(t, 1, resp.TotalPages)
			assert.EqualValues(t, 6, resp.TotalRecordCount)
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
			[]string{
				dbtools.FixtureDellR640BMC.ID,
				dbtools.FixtureDellR640BIOS.ID,
				dbtools.FixtureDellR6515BMC.ID,
				dbtools.FixtureDellR6515BIOS.ID,
				dbtools.FixtureDellR640CPLD.ID,
			},
			false,
			"",
		},
		{
			"search by model",
			&serverservice.ComponentFirmwareVersionListParams{
				Model: []string{"X11DPH-T"},
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
			[]string{dbtools.FixtureDellR6515BIOS.ID},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, _, err := s.Client.ListServerComponentFirmware(context.TODO(), tt.params)
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
		fw, _, err := s.Client.GetServerComponentFirmware(ctx, uuid.MustParse(dbtools.FixtureDellR640BMC.ID))

		if !expectError {
			require.NoError(t, err)
			assert.Equal(t, fw.UUID, uuid.MustParse(dbtools.FixtureDellR640BMC.ID))
		}

		return err
	})
}

func TestIntegrationServerComponentFirmwareCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		testFirmware := serverservice.ComponentFirmwareVersion{
			UUID:          uuid.New(),
			Vendor:        "dell",
			Model:         []string{"r615"},
			Filename:      "foobar",
			Version:       "21.07.00",
			Component:     "system",
			Checksum:      "foobar",
			UpstreamURL:   "https://vendor.com/firmwares/DSU_21.07.00/",
			RepositoryURL: "http://example-firmware-bucket.s3.amazonaws.com/firmware/dell/DSU_21.07.00/",
		}

		id, resp, err := s.Client.CreateServerComponentFirmware(ctx, testFirmware)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, id)
			assert.Equal(t, testFirmware.UUID.String(), id.String())
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/server-component-firmwares/%s", id), resp.Links.Self.Href)
		}

		return err
	})

	var testCases = []struct {
		testName         string
		firmware         *serverservice.ComponentFirmwareVersion
		expectedError    bool
		expectedResponse string
		errorMsg         string
	}{
		{
			"empty required parameters",
			&serverservice.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "dell",
				Model:         nil,
				Filename:      "foobar",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
			},
			true,
			"400",
			"Error:Field validation for 'Model' failed on the 'required' tag",
		},
		{
			"required lowercase parameters",
			&serverservice.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "DELL",
				Model:         []string{"r615"},
				Filename:      "foobar",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
			},
			true,
			"400",
			"Error:Field validation for 'Vendor' failed on the 'lowercase' tag",
		},
		{
			"filename allowed to be mixed case",
			&serverservice.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "dell",
				Model:         []string{"r615"},
				Filename:      "fooBAR",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
			},
			false,
			"200",
			"",
		},
		{
			"duplicate vendor/model/version not allowed",
			&serverservice.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "dell",
				Model:         []string{"r615"},
				Filename:      "foobar",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
			},
			true,
			"400",
			"unable to insert into component_firmware_version: pq: duplicate key value violates unique constraint \"vendor_model_version_unique\"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			fwUUID, _, err := s.Client.CreateServerComponentFirmware(context.TODO(), *tt.firmware)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Contains(t, err.Error(), tt.expectedResponse)
				return
			}
			assert.Equal(t, tt.firmware.UUID.String(), fwUUID.String())
		})
	}
}

func TestIntegrationServerComponentFirmwareDelete(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)
		_, err := s.Client.DeleteServerComponentFirmware(ctx, serverservice.ComponentFirmwareVersion{UUID: uuid.MustParse(dbtools.FixtureDellR640BMC.ID)})

		return err
	})
}

func TestIntegrationServerComponentFirmwareUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		fw := serverservice.ComponentFirmwareVersion{
			UUID:          uuid.MustParse(dbtools.FixtureDellR640BMC.ID),
			Vendor:        "dell",
			Model:         []string{"r615"},
			Filename:      "foobarino",
			Version:       "21.07.00",
			Component:     "bios",
			Checksum:      "foobar",
			UpstreamURL:   "https://vendor.com/firmware-file",
			RepositoryURL: "https://example-firmware-bucket.s3.amazonaws.com/firmware/dell/r615/bios/filename.ext",
		}

		resp, err := s.Client.UpdateServerComponentFirmware(ctx, uuid.MustParse(dbtools.FixtureDellR640BMC.ID), fw)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/server-component-firmwares/%s", dbtools.FixtureDellR640BMC.ID), resp.Links.Self.Href)
			fw, _, _ := s.Client.GetServerComponentFirmware(ctx, uuid.MustParse(dbtools.FixtureDellR640BMC.ID))
			assert.Equal(t, "foobarino", fw.Filename)
		}

		return err
	})
}
