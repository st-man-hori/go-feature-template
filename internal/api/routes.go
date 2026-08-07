package api

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/st-man-hori/go-feature-template/internal/features/task"
)

// RegisterRoutes は各 Feature のエンドポイントを登録する。
//
// Laravel比較: routes/api.php の中身そのものに相当する。Laravel版がこのファイルに
// 全 Feature のルートを Route::prefix() ブロックで並べているのと同じく、ここに
// Feature ごとのブロックを並べる。Feature を追加するときは、このファイルに
// ブロックを1つ足すだけでよい。
func RegisterRoutes(r chi.Router, db *gorm.DB) {
	// --- Task (サンプル Feature) -------------------------------------------
	taskHandler := task.NewHandler(task.NewUseCase(db))
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", taskHandler.Index)
		r.Post("/", taskHandler.Store)
		r.Get("/{task}", taskHandler.Show)
		r.Patch("/{task}", taskHandler.Update)
		r.Delete("/{task}", taskHandler.Destroy)
	})
}
