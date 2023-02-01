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
	_ "github.com/lib/pq" // Register the Postgres driver.
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gocloud.dev/secrets"

	// import gocdk secret drivers
	_ "gocloud.dev/secrets/localsecrets"

	"go.hollow.sh/serverservice/internal/models"
)

// TestDBURI is the URI for the test database
var TestDBURI = os.Getenv("SERVERSERVICE_CRDB_URI")
var testDB *sqlx.DB
var testKeeper *secrets.Keeper

func testDatastore(t *testing.T) error {
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

	return addFixtures(t)
}

// TestSecretKeeper will return the secret keeper we are using for this test run. This allows
// use to use the same one for the entire test run so the secrets are able to be decrypted.
func TestSecretKeeper(t *testing.T) *secrets.Keeper {
	if testKeeper != nil {
		return testKeeper
	}

	keeper, err := secrets.OpenKeeper(context.TODO(), "base64key://")
	require.NoError(t, err)

	testKeeper = keeper

	return keeper
}

// DatabaseTest allows you to run tests that interact with the database
func DatabaseTest(t *testing.T) *sqlx.DB {
	RegisterHooks()

	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	t.Cleanup(func() {
		cleanDB()
		err := addFixtures(t)
		require.NoError(t, err, "Unexpected error setting up fixture data")
	})

	err := testDatastore(t)
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
	models.ServerCredentials().DeleteAll(ctx, testDB)
	models.Servers(qm.WithDeleted()).DeleteAll(ctx, testDB, true)
	models.ComponentFirmwareVersions().DeleteAll(ctx, testDB)
	models.ComponentFirmwareSets().DeleteAll(ctx, testDB)
	models.ComponentFirmwareSetMaps().DeleteAll(ctx, testDB)
	// don't delete the builtin ServerCredentialTypes. Those are expected to exist for the application to work
	models.ServerCredentialTypes(models.ServerCredentialTypeWhere.Builtin.EQ(false)).DeleteAll(ctx, testDB)
	models.ServerConditionTypes().DeleteAll(ctx, testDB)
	models.ServerConditionStatusTypes().DeleteAll(ctx, testDB)
	models.ServerConditions().DeleteAll(ctx, testDB)
	testDB.Exec("SET sql_safe_updates = true;")
}
