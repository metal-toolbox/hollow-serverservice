package serverservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.hollow.sh/serverservice/internal/dbtools"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

func TestIntegrationServerConditionList(t *testing.T) {
	s := serverTest(t)

	serverUUID := uuid.MustParse(dbtools.FixtureNemo.ID)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		r, resp, err := s.Client.ListServerConditions(ctx, serverUUID, nil)
		if !expectError {
			assert.NoError(t, err)
			assert.Len(t, r, 1)

			assert.EqualValues(t, 1, resp.PageCount)
			assert.EqualValues(t, 1, resp.TotalPages)
			assert.EqualValues(t, 1, resp.TotalRecordCount)
			// We returned everything, so we shouldn't have a next page info
			assert.Nil(t, resp.Links.Next)
			assert.Nil(t, resp.Links.Previous)
		}

		return err
	})
}

func TestIntegrationServerConditionGet(t *testing.T) {
	s := serverTest(t)

	serverUUID := uuid.MustParse(dbtools.FixtureNemo.ID)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		r, resp, err := s.Client.GetServerCondition(ctx, serverUUID, dbtools.FixtureSwimConditionType.Slug)
		if !expectError {
			assert.NoError(t, err)
			assert.Equal(t, r.Slug, dbtools.FixtureSwimConditionType.Slug)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/servers/%s/conditions/%s", serverUUID, r.Slug), resp.Links.Self.Href)
			assert.Nil(t, resp.Links.Next)
			assert.Nil(t, resp.Links.Previous)
		}

		return err
	})
}

func TestIntegrationServerConditionDelete(t *testing.T) {
	s := serverTest(t)

	serverUUID := uuid.MustParse(dbtools.FixtureNemo.ID)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		resp, err := s.Client.DeleteServerCondition(ctx, serverUUID, dbtools.FixtureSwimConditionType.Slug)
		if !expectError {
			assert.NoError(t, err)
			assert.Equal(t, resp.Message, "resource deleted")
		}

		return err
	})
}

func TestIntegrationServerConditionSet(t *testing.T) {
	s := serverTest(t)

	serverUUID := uuid.MustParse(dbtools.FixtureNemo.ID)

	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		condition := &serverservice.ServerCondition{
			Slug:       dbtools.FixtureSwimConditionType.Slug,
			Status:     dbtools.FixtureSwimConditionStatusTypeSucceeded.Slug,
			Parameters: []byte(`{"distance_km":"5"}`),
		}

		resp, err := s.Client.SetServerCondition(ctx, serverUUID, condition)
		if !expectError {
			assert.NoError(t, err)
			assert.Equal(t, "resource updated", resp.Message)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/servers/%s/conditions/%s", serverUUID, dbtools.FixtureSwimConditionType.Slug), resp.Links.Self.Href)
		}

		return err
	})

	cases := []struct {
		name        string
		serverID    uuid.UUID
		condition   serverservice.ServerCondition
		errContains string
	}{
		{
			"create",
			uuid.MustParse(dbtools.FixtureMarlin.ID),
			serverservice.ServerCondition{
				Slug:       dbtools.FixtureSwimConditionType.Slug,
				Status:     dbtools.FixtureSwimConditionStatusTypeActive.Slug,
				Parameters: []byte(`{"foo":"bar"}`),
			},
			"",
		},

		{
			"update",
			uuid.MustParse(dbtools.FixtureMarlin.ID),
			serverservice.ServerCondition{
				Slug:         dbtools.FixtureSwimConditionType.Slug,
				Status:       dbtools.FixtureSwimConditionStatusTypeSucceeded.Slug,
				StatusOutput: []byte(`{"all":"done"}`),
			},
			"",
		},
		{
			"create fails on unknown server",
			uuid.New(),
			serverservice.ServerCondition{
				Slug:       dbtools.FixtureSwimConditionType.Slug,
				Status:     dbtools.FixtureSwimConditionStatusTypeActive.Slug,
				Parameters: []byte(`{"foo":"bar"}`),
			},
			"server not found",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.Client.SetServerCondition(
				context.TODO(),
				tc.serverID,
				&tc.condition,
			)
			if err != nil {
				assert.ErrorContains(t, err, tc.errContains)
				return
			}

			got, _, err := s.Client.GetServerCondition(context.TODO(), tc.serverID, dbtools.FixtureSwimConditionType.Slug)
			assert.NoError(t, err)

			assert.Equal(t, tc.condition.Slug, got.Slug)
			assert.Equal(t, tc.condition.Status, got.Status)

			if len(tc.condition.Parameters) > 0 {
				assert.Equal(t, tc.condition.Parameters, got.Parameters)
			}

			if len(tc.condition.StatusOutput) > 0 {
				assert.Equal(t, tc.condition.StatusOutput, got.StatusOutput)
			}
		})
	}
}
