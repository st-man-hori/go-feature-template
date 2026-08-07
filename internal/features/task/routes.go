package task

import "github.com/go-chi/chi/v5"

// RegisterRoutes mounts the Task feature's endpoints under whatever prefix
// the caller passes in — the equivalent of the Route::prefix('tasks')
// block in Laravel's routes/api.php, except it lives inside the Feature
// directory itself instead of a shared routes file, so nothing about Task
// needs to be found by looking outside internal/features/task/.
func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", h.Index)
		r.Post("/", h.Store)
		r.Get("/{task}", h.Show)
		r.Patch("/{task}", h.Update)
		r.Delete("/{task}", h.Destroy)
	})
}
