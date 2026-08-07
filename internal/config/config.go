// Package config は環境変数からアプリケーション設定を読み込む。
//
// Laravel比較: config/database.php + env() ヘルパーによる .env 読み込みを置き換える。
// 設定キャッシュのステップ(php artisan config:cache)は存在しない。Go はプロセス
// 起動時に環境変数を一度だけ読むだけで、リクエストごとの起動コストをキャッシュで
// 避けるという概念自体が無いため。
package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv string
	Port   string

	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string
}

func Load() Config {
	return Config{
		AppEnv: getEnv("APP_ENV", "local"),
		Port:   getEnv("APP_PORT", "8000"),

		DBHost:     getEnv("DB_HOST", "mysql"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBDatabase: getEnv("DB_DATABASE", "app"),
		DBUsername: getEnv("DB_USERNAME", "app"),
		DBPassword: getEnv("DB_PASSWORD", "secret"),
	}
}

// DSN は gorm.io/driver/mysql 用の MySQL DSN(データソース名)を返す。
func (c Config) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUsername, c.DBPassword, c.DBHost, c.DBPort, c.DBDatabase,
	)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
