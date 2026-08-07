package task

import "github.com/go-chi/chi/v5"

// RegisterRoutes は、呼び出し側が渡したプレフィックスの下に Task Feature の
// エンドポイントをマウントする — Laravel の routes/api.php にある
// Route::prefix('tasks') ブロックに相当するが、共有のルートファイルではなく
// Feature ディレクトリ自体の中に置かれている点が異なる。そのため、Task に
// 関することを探すのに internal/features/task/ の外を見る必要がない。
func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", h.Index)
		r.Post("/", h.Store)
		r.Get("/{task}", h.Show)
		r.Patch("/{task}", h.Update)
		r.Delete("/{task}", h.Destroy)
	})
}
