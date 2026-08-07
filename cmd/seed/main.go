// Command seed populates the database with sample data.
//
// Laravel比較: `php artisan db:seed` / database/seeders/DatabaseSeeder.php に
// 相当する。あちらが本物のフィクスチャではなく薄いプレースホルダーとして
// 存在するのと同じく、こちらもサンプルタスクを数件挿入するだけ — 実際の
// Feature ができたら置き換えること。
package main

import (
	"log/slog"
	"os"

	"github.com/st-man-hori/go-feature-template/internal/config"
	"github.com/st-man-hori/go-feature-template/internal/database"
	"github.com/st-man-hori/go-feature-template/internal/models"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	tasks := []models.Task{
		{Title: "Write the quarterly report"},
		{Title: "Review pull requests", IsDone: true},
	}

	if err := db.Create(&tasks).Error; err != nil {
		slog.Error("failed to seed database", "error", err)
		os.Exit(1)
	}

	slog.Info("seeded database", "count", len(tasks))
}
