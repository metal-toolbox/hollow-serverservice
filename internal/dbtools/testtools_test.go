package dbtools_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"
)

// TestDatabaseTest is used to force the test commands in this package to run during
// it's test cycle. It's also used to to some basic checks to make sure everything
// for our test database setup is correct.
func TestDatabaseTest(t *testing.T) {
	ctx := context.TODO()

	t.Run("make changes to Nemo", func(t *testing.T) {
		db := dbtools.DatabaseTest(t)
		dbtools.FixtureNemo.Name = null.StringFrom("Something New")

		_, err := dbtools.FixtureNemo.Update(ctx, db, boil.Infer())
		assert.NoError(t, err)

		srv, err := models.FindServer(ctx, db, dbtools.FixtureNemo.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Something New", srv.Name.String)
	})

	t.Run("new test should be restored", func(t *testing.T) {
		db := dbtools.DatabaseTest(t)
		srv, err := models.FindServer(ctx, db, dbtools.FixtureNemo.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Nemo", srv.Name.String)
	})
}
