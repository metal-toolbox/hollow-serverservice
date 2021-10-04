//go:build testtools
// +build testtools

package dbtools

import (
	"context"
	"os"
	"testing"

	// import the crdbpgx for automatic retries of errors for crdb that support retry
	_ "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.hollow.sh/serverservice/internal/models"
)

// TestDBURI is the URI for the test database
var TestDBURI = os.Getenv("SERVERSERVICE_DB_URI")
var testDB *sqlx.DB

func testDatastore() error {
	// don't setup the datastore if we already have one
	if testDB != nil {
		return nil
	}

	// Uncomment when you are having database issues with your tests and need to see the db logs
	// Hidden by default because it can be noisy and make it harder to read normal failures.
	// You can also enable at the beginning of your test and then disable it again at the end
	// boil.DebugMode = true

	db, err := sqlx.Open("postgres", TestDBURI)
	if err != nil {
		return err
	}

	testDB = db

	cleanDB()

	return addFixtures()
}

// DatabaseTest allows you to run tests that interact with the database
func DatabaseTest(t *testing.T) *sqlx.DB {
	RegisterHooks()

	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	t.Cleanup(func() {
		cleanDB()
		err := addFixtures()
		require.NoError(t, err, "Unexpected error setting up fixture data")
	})

	err := testDatastore()
	require.NoError(t, err, "Unexpected error getting connection to test datastore")

	return testDB
}

// nolint
func cleanDB() {
	ctx := context.TODO()
	// Make sure the deletion goes in order so you don't break the databases foreign key constraints
	testDB.Exec("SET sql_safe_updates = false;")
	models.Attributes().DeleteAll(ctx, testDB)
	models.VersionedAttributes().DeleteAll(ctx, testDB)
	models.ServerComponents().DeleteAll(ctx, testDB)
	models.ServerComponentTypes().DeleteAll(ctx, testDB)
	models.Servers(qm.WithDeleted()).DeleteAll(ctx, testDB, true)
	testDB.Exec("SET sql_safe_updates = true;")
}
