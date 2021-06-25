package db

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

// NewTestStore gives us a database instance for testing. We return the database
// so that we can manually manage transactions for tests to rollback changes to
// enable starting with a clean database for each tests
func NewTestStore(uri string, lg *zap.Logger) (*gorm.DB, error) {
	var err error

	zl := zapgorm2.New(lg)
	zl.SetAsDefault()

	pg := postgres.New(postgres.Config{
		DSN:                  uri,
		PreferSimpleProtocol: true,
	})

	db, err = gorm.Open(pg, &gorm.Config{
		Logger: zl,
	})
	if err != nil {
		return nil, err
	}

	if err := db.Exec("CREATE DATABASE IF NOT EXISTS hollow_test;").Error; err != nil {
		fmt.Println("test database failed to create")
	}

	err = migrate()

	return db, err
}
