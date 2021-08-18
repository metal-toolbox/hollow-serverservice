//+build testtools

package dbtools

import (
	"context"
	"database/sql"
	"os"
	"testing"

	// pq is imported to get our database connection
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"go.metalkube.net/hollow/internal/models"
)

// TestDBURI is the URI for the test database
var TestDBURI = os.Getenv("HOLLOW_TEST_DB")
var testDB *sql.DB

func testDatastore() error {
	// don't setup the datastore if we already have one
	if testDB != nil {
		return nil
	}

	// Uncomment when you are having database issues with your tests and need to see the db logs
	// Hidden by default because it can be noisy and make it harder to read normal failures.
	// You can also enable at the beginning of your test and then disable it again at the end
	// boil.DebugMode = true

	db, err := sql.Open("postgres", TestDBURI)
	if err != nil {
		return err
	}

	testDB = db

	cleanDB()

	return addFixtures()
}

// DatabaseTest allows you to run tests that interact with the database
func DatabaseTest(t *testing.T) *sql.DB {
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
	models.Servers().DeleteAll(ctx, testDB)
	testDB.Exec("SET sql_safe_updates = true;")
}
