package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"moul.io/zapgorm2"
)

var db *gorm.DB

// Tabler provides an interface for GORM that allows use to change the name a table
// uses. For example by default Hardware gets stored in a table called "hardwares"
// but we want to store it in one called "hardware".
type Tabler interface {
	TableName() string
}

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

// Ping checks to ensure that the database is available and processing queries
func Ping() bool {
	res := db.Exec("select 1 as result;")
	if res.Error != nil {
		return false
	}

	return res.RowsAffected == int64(1)
}

func migrate() error {
	return db.AutoMigrate(
		&Attributes{},
		&BIOSConfig{},
		&Hardware{},
		&HardwareComponent{},
		&HardwareComponentType{},
	)
}
