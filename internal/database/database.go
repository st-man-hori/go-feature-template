// Package database opens the application's *gorm.DB connection.
//
// Laravel comparison: replaces config/database.php + the DB manager Laravel
// wires up automatically. Schema is owned by migrations/*.sql (run via
// `make migrate`, see the Laravel template's database/migrations/), not by
// GORM's AutoMigrate — AutoMigrate is only used in tests (see internal/testutil)
// against an in-memory SQLite DB, similar to how the Laravel template runs
// PHPUnit against SQLite instead of the real MySQL service.
package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/st-man-hori/go-feature-template/internal/config"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	logLevel := logger.Warn
	if cfg.AppEnv == "local" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return db, nil
}
