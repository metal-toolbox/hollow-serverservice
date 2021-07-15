//+build testtools

package db

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TestDBURI is the URI for the test database
var TestDBURI = "postgresql://root@localhost:26257/hollow_test?sslmode=disable"

func testDatastore() error {
	// don't setup the datastore if we already have one
	if db != nil {
		return nil
	}

	err := NewPostgresStore(TestDBURI, zap.NewNop())
	if err != nil {
		return err
	}

	cleanDB()

	return setupTestData()
}

// DatabaseTest allows you to run tests that interact with the database
func DatabaseTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	t.Cleanup(func() {
		cleanDB()
		err := setupTestData()
		require.NoError(t, err, "Unexpected error setting up test data")
	})

	err := testDatastore()
	require.NoError(t, err, "Unexpected error setting up test datastore")
}

func cleanDB() {
	d := db.Session(&gorm.Session{AllowGlobalUpdate: true})
	// Make sure the deletion goes in order so you don't break the databases foreign key constraints
	d.Delete(&Attributes{})
	d.Delete(&BIOSConfig{})
	d.Delete(&HardwareComponent{})
	d.Delete(&HardwareComponentType{})
	d.Delete(&Hardware{})
}
