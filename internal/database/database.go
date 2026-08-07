// Package database はアプリケーションの *gorm.DB 接続を開く。
//
// Laravel比較: config/database.php + Laravel が自動で組み立てる DB マネージャーを
// 置き換える。スキーマは migrations/*.sql(`make migrate` で実行、Laravel版の
// database/migrations/ に相当)が管理しており、GORM の AutoMigrate では管理しない
// — AutoMigrate はテスト(internal/testutil 参照)で in-memory SQLite に対して
// のみ使う。Laravel版が PHPUnit を本物の MySQL ではなく SQLite に対して実行する
// のと同じ理由。
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
