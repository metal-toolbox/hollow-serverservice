package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"moul.io/zapgorm2"
)

var db *gorm.DB

// NewPostgresStore creates a new PostgeSQL store instance, opening a connection and
// applying any db migrations available.
func NewPostgresStore(uri string, lg *zap.Logger) error {
	var err error

	zl := zapgorm2.New(lg)
	zl.SetAsDefault()

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  uri,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: zl,
	})
	if err != nil {
		return err
	}

	if err := db.Use(prometheus.New(prometheus.Config{})); err != nil {
		return err
	}

	return migrate()
}

func migrate() error {
	return db.AutoMigrate(
		&BIOSConfig{},
		&Hardware{},
	)
}
