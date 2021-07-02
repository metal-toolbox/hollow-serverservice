package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.metalkube.net/hollow/internal/db"
)

func TestNewPostgresStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	err := db.NewPostgresStore(db.TestDBURI, zap.NewNop())
	assert.NoError(t, err)
}

func TestNewPostgresStoreFailure(t *testing.T) {
	err := db.NewPostgresStore("invalid-uri", zap.NewNop())
	assert.Error(t, err)
}

func TestPing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	var testCases = []struct {
		testName       string
		dbURI          string
		expectedResult bool
	}{
		{"happy path", db.TestDBURI, true},
	}

	for _, tt := range testCases {
		err := db.NewPostgresStore(tt.dbURI, zap.NewNop())
		require.NoError(t, err)

		res := db.Ping()
		assert.Equal(t, tt.expectedResult, res, tt.testName)
	}
}
