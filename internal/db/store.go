package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"moul.io/zapgorm2"
)

// Tabler provides an interface for GORM that allows use to change the name a table
// uses. For example by default Hardware gets stored in a table called "hardwares"
// but we want to store it in one called "hardware".
type Tabler interface {
	TableName() string
}

// Store provides access to the underlaying datastore
type Store struct {
	db *gorm.DB
}

// NewPostgresStore creates a new PostgeSQL store instance, opening a connection and
// applying any db migrations available.
func NewPostgresStore(uri string, lg *zap.Logger) (*Store, error) {
	zl := zapgorm2.New(lg)
	zl.SetAsDefault()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  uri,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: zl,
	})
	if err != nil {
		return nil, err
	}

	if err := db.Use(prometheus.New(prometheus.Config{})); err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

// Ping checks to ensure that the database is available and processing queries
func (s *Store) Ping() bool {
	if s.db == nil {
		return false
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}
