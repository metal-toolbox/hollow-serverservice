package serverservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	serverservice "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

func TestIntegrationServerCredentialsUpsert(t *testing.T) {
	ctx := context.TODO()
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		id := uuid.MustParse(dbtools.FixtureDory.ID)
		slug := serverservice.ServerCredentialTypeBMC
		path := fmt.Sprintf("http://test.hollow.com/api/v1/servers/%s/credentials/%s", id, slug)

		resp, err := s.Client.SetCredential(ctx, id, slug, "postgre", "something")
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, path, resp.Links.Self.Href)
		}

		return err
	})

	s.Client.SetToken(validToken(adminScopes))

	t.Run("inserts when secret doesn't exist", func(t *testing.T) {
		uuid := uuid.MustParse(dbtools.FixtureMarlin.ID)
		slug := serverservice.ServerCredentialTypeBMC

		// Make sure our server doesn't already have a BMC secret
		secret, _, err := s.Client.GetCredential(ctx, uuid, slug)
		require.Error(t, err)
		require.Nil(t, secret)
		require.Contains(t, err.Error(), "not found")

		// Create the secret
		_, err = s.Client.SetCredential(ctx, uuid, slug, "postgre", "supersecret")
		assert.NoError(t, err)

		// Ensure we can retrieve the secret we set
		secret, _, err = s.Client.GetCredential(ctx, uuid, slug)
		assert.NoError(t, err)
		assert.Equal(t, "supersecret", secret.Password)
		assert.Equal(t, "postgre", secret.Username)
	})

	t.Run("updates when secret already exist", func(t *testing.T) {
		uuid := uuid.MustParse(dbtools.FixtureNemo.ID)
		slug := serverservice.ServerCredentialTypeBMC

		// Get the existing secret
		secret, _, err := s.Client.GetCredential(ctx, uuid, slug)
		assert.NoError(t, err)
		assert.NotEqual(t, "mynewSecret!", secret.Password)

		// Update the secret
		_, err = s.Client.SetCredential(ctx, uuid, slug, "foobar", "mynewSecret!")
		assert.NoError(t, err)

		// Get the new secret
		newSecret, _, err := s.Client.GetCredential(ctx, uuid, slug)
		assert.NoError(t, err)
		assert.Equal(t, "mynewSecret!", newSecret.Password)
		assert.Equal(t, "foobar", newSecret.Username)
		// ensure timestamps were updated correctly
		assert.Equal(t, secret.CreatedAt, newSecret.CreatedAt)
		assert.NotEqual(t, newSecret.UpdatedAt, secret.UpdatedAt)
		assert.Greater(t, newSecret.UpdatedAt, secret.UpdatedAt)
	})

	t.Run("fails if server uuid not found", func(t *testing.T) {
		uuid := uuid.New()
		slug := serverservice.ServerCredentialTypeBMC

		// Make sure our server doesn't already have a BMC secret
		_, err := s.Client.SetCredential(ctx, uuid, slug, "postgre", "secret")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "server not found")
	})

	t.Run("fails if secret type slug not found", func(t *testing.T) {
		uuid := uuid.MustParse(dbtools.FixtureMarlin.ID)
		slug := "notfound"

		_, err := s.Client.SetCredential(ctx, uuid, slug, "postgre", "secret")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestIntegrationServerSecretsDelete(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		id := uuid.MustParse(dbtools.FixtureNemo.ID)
		slug := serverservice.ServerCredentialTypeBMC

		_, err := s.Client.DeleteCredential(ctx, id, slug)
		if !expectError {
			assert.NoError(t, err)
		}

		return err
	})
}
