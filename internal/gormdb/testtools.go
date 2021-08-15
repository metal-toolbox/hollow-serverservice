//+build testtools

package gormdb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TestDBURI is the URI for the test database
var TestDBURI = os.Getenv("HOLLOW_TEST_DB")
var testStore *Store

func testDatastore() error {
	// don't setup the datastore if we already have one
	if testStore != nil {
		return nil
	}

	// Uncomment when you are having database issues with your tests and need to see the db logs
	// Hidden by default because it can be noisy and make it harder to read normal failures
	// l, err := zap.NewDevelopment()
	// if err != nil {
	// 	return err
	// }

	l := zap.NewNop()

	s, err := NewPostgresStore(TestDBURI, l)
	if err != nil {
		return err
	}

	testStore = s

	cleanDB()

	return s.setupTestData()
}

// DatabaseTest allows you to run tests that interact with the database
func DatabaseTest(t *testing.T) *Store {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	t.Cleanup(func() {
		cleanDB()
		err := testStore.setupTestData()
		require.NoError(t, err, "Unexpected error setting up test data")
	})

	err := testDatastore()
	require.NoError(t, err, "Unexpected error setting up test datastore")

	return testStore
}

func cleanDB() {
	d := testStore.db.Session(&gorm.Session{AllowGlobalUpdate: true})
	// Make sure the deletion goes in order so you don't break the databases foreign key constraints
	d.Delete(&Attributes{})
	d.Delete(&VersionedAttributes{})
	d.Delete(&ServerComponent{})
	d.Delete(&ServerComponentType{})
	d.Unscoped().Delete(&Server{})
}

// DatabaseTestErr provides a test database
func DatabaseTestErr() (*Store, error) {
	if testing.Short() {
		return nil, nil
	}

	if err := testDatastore(); err != nil {
		return nil, err
	}

	return testStore, nil
}
