package hollow_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.metalkube.net/hollow/internal/db"
	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

func TestIntegrationServerCreateAttributes(t *testing.T) {
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

func TestIntegrationServerListAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		attrs, resp, err := s.Client.Server.ListAttributes(ctx, db.FixtureServerNemo.ID)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, resp.Links.Self)
			assert.Len(t, attrs, 2)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/servers/%s/attributes", db.FixtureServerNemo.ID), resp.Links.Self.Href)
		}

		return err
	})

	var testCases = []struct {
		testName string
		srvUUID  uuid.UUID
		errorMsg string
	}{
		{
			"returns not found on missing server uuid",
			uuid.New(),
			"response code: 404",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			_, _, err := s.Client.Server.ListAttributes(context.TODO(), tt.srvUUID)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestIntegrationServerGetAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		attr, resp, err := s.Client.Server.GetAttributes(ctx, db.FixtureServerNemo.ID, db.FixtureNamespaceMetadata)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, resp.Links.Self)
			assert.ElementsMatch(t, attr.Data, db.FixtureAttributesNemoMetadata.Data)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/servers/%s/attributes/%s", db.FixtureServerNemo.ID, db.FixtureNamespaceMetadata), resp.Links.Self.Href)
		}

		return err
	})
}

func TestIntegrationServerDeleteAttributes(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		_, err := s.Client.Server.DeleteAttributes(ctx, db.FixtureServerNemo.ID, db.FixtureNamespaceMetadata)
		if !expectError {
			require.NoError(t, err)

			// ensure the attributes are gone
			_, _, err2 := s.Client.Server.GetAttributes(ctx, db.FixtureServerNemo.ID, db.FixtureNamespaceMetadata)
			require.Error(t, err2)
			assert.Contains(t, err2.Error(), "response code: 404, message: resource not found")
		}

		return err
	})
}

func TestIntegrationServerUpdateAttributes(t *testing.T) {
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
