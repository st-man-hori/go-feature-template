// Package testutil は、in-memory SQLite データベースに対して HTTP 層を
// 立ち上げるテスト専用のヘルパーを提供する。
//
// Laravel比較: Laravel版の phpunit.xml が RefreshDatabase を本物の MySQL では
// なく SQLite に向けているのと同じ狙いの Go 版で、Feature テストを Docker 無しで
// 高速に実行できる。本番のスキーマは migrations/*.sql(MySQL 専用の SQL)が
// 管理しているため、テストはそれらのマイグレーションを別の SQL 方言で再生する
// のではなく、モデル構造体からの GORM の AutoMigrate を使う。
package testutil

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/st-man-hori/go-feature-template/internal/api"
	"github.com/st-man-hori/go-feature-template/internal/models"
)

// NewDB は1テストごとに独立した in-memory SQLite データベースを新規に開き、
// GORM のモデル構造体からマイグレーションする。SQLite の :memory: モードは
// 接続ごとにまっさらな新しいデータベースを作るため、コネクションプールを
// 1本に固定している — そうしないと、テスト途中で開かれた2本目のコネクションが
// 空の未マイグレーション状態のデータベースを見てしまう。
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

// NewRouter は internal/api.RegisterRoutes をそのまま使って /api 配下の
// 全ルートを持つ chi.Mux を db から組み立てる。internal/api.NewRouter と
// ルート表を共有しつつ、ログや Swagger のミドルウェアは含めない —
// テスト出力のノイズになるだけなので。
func NewRouter(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		api.RegisterRoutes(r, db)
	})

	return r
}
