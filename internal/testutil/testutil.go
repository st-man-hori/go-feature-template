// Package testutil provides test-only helpers for spinning up the HTTP
// layer against an in-memory SQLite database.
//
// Laravel comparison: this is the Go equivalent of the Laravel template's
// phpunit.xml pointing RefreshDatabase at SQLite instead of the real MySQL
// service, so feature tests run fast with no Docker dependency. Production
// schema is owned by migrations/*.sql (MySQL-specific SQL), so tests use
// GORM's AutoMigrate against the model structs instead of replaying those
// migrations under a different SQL dialect.
package testutil

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/st-man-hori/go-feature-template/internal/features/task"
	"github.com/st-man-hori/go-feature-template/internal/models"
)

// NewDB opens a fresh, isolated in-memory SQLite database for one test and
// migrates it from the GORM model structs. SQLite's :memory: mode creates a
// brand-new empty database per connection, so the pool is pinned to a
// single connection — otherwise a second connection opened mid-test would
// see an empty, unmigrated database.
func NewDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(&models.Task{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	return db
}

// NewRouter builds a chi.Mux with every route the app mounts under /api,
// backed by db. It mirrors internal/api.NewRouter's route table without
// that package's logging/Swagger middleware, which would just add noise to
// test output.
func NewRouter(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		h := task.NewHandler(task.NewUseCase(db))
		task.RegisterRoutes(r, h)
	})

	return r
}
