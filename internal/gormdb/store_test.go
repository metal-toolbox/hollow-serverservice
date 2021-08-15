package gormdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.metalkube.net/hollow/internal/gormdb"
)

func TestNewPostgresStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	s, err := gormdb.NewPostgresStore(gormdb.TestDBURI, zap.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestNewPostgresStoreFailure(t *testing.T) {
	s, err := gormdb.NewPostgresStore("invalid-uri", zap.NewNop())
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestPingNoDB(t *testing.T) {
	var testCases = []struct {
		testName       string
		expectedResult bool
	}{
		{"no db configured, return false", false},
	}

	for _, tt := range testCases {
		s := &gormdb.Store{}
		res := s.Ping()
		assert.Equal(t, tt.expectedResult, res, tt.testName)
	}
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
		{"happy path", gormdb.TestDBURI, true},
	}

	for _, tt := range testCases {
		s, err := gormdb.NewPostgresStore(tt.dbURI, zap.NewNop())
		require.NoError(t, err)
		require.NotNil(t, s)

		res := s.Ping()
		assert.Equal(t, tt.expectedResult, res, tt.testName)
	}
}
