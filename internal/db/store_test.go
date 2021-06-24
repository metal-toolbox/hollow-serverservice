package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.metalkube.net/hollow/internal/db"
)

var testDB *gorm.DB

var testDBURI = "postgresql://root@localhost:26257/hollow_test?sslmode=disable"

func testDatastore() {
	var err error

	// don't setup the datastore if we already have one
	if testDB != nil {
		return
	}

	testDB, err = db.NewTestStore(testDBURI, zap.NewNop())
	if err != nil {
		panic(err)
	}

	cleanDB()

	if err := setupTestData(); err != nil {
		panic(err)
	}
}

func databaseTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	t.Cleanup(func() {
		cleanDB()
		err := setupTestData()
		assert.NoError(t, err, "Unexpected error setting up test data")
	})

	testDatastore()
}

func cleanDB() {
	d := testDB.Session(&gorm.Session{AllowGlobalUpdate: true})
	d.Delete(&db.Hardware{})
}

func TestNewPostgresStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	err := db.NewPostgresStore(testDBURI, zap.NewNop())
	assert.NoError(t, err)
}

func TestNewPostgresStoreFailure(t *testing.T) {
	err := db.NewPostgresStore("invalid-uri", zap.NewNop())
	assert.Error(t, err)
}
