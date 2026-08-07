// Package config loads application configuration from environment variables.
//
// Laravel comparison: this replaces config/database.php + .env reading via
// the env() helper. There is no config caching step (php artisan config:cache)
// because Go reads env vars once at process startup — there's no per-request
// bootstrap cost to cache away.
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

// DSN returns a MySQL data source name suitable for gorm.io/driver/mysql.
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
