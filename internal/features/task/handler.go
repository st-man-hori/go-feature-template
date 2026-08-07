package task

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/st-man-hori/go-feature-template/internal/apperror"
	"github.com/st-man-hori/go-feature-template/internal/httpx"
)

// Handler is the equivalent of Laravel's TaskController: a thin layer that
// decodes/validates the request, calls the UseCase, and shapes the
// response. One method = one action, same as the Laravel template.
type Handler struct {
	uc *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{uc: uc}
}

// Index lists tasks, optionally filtered by completion status.
//
//	@Summary	List tasks
//	@Tags		tasks
//	@Produce	json
//	@Param		isDone	query		bool	false	"Filter by completion status"
//	@Param		perPage	query		int		false	"Results per page (1-100, default 15)"
//	@Success	200		{object}	TaskListDataResponse
//	@Failure	422		{object}	httpx.ValidationErrorResponse
//	@Router		/tasks [get]
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	req, err := NewIndexTaskRequest(r.URL.Query())
	if err != nil {
		httpx.Error(w, err)
		return
	}

	tasks, err := h.uc.Index(req)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.Data(w, http.StatusOK, toTaskResponses(tasks))
}

// Show returns a single task by id.
//
//	@Summary	Get a task
//	@Tags		tasks
//	@Produce	json
//	@Param		task	path		int	true	"Task ID"
//	@Success	200		{object}	TaskDataResponse
//	@Failure	404		{object}	httpx.ErrorResponse
//	@Router		/tasks/{task} [get]
func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	id, err := taskIDFromPath(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	task, err := h.uc.Show(id)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.Data(w, http.StatusOK, toTaskResponse(task))
}

// Store creates a new task.
//
//	@Summary	Create a task
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		task	body		StoreTaskRequest	true	"Task to create"
//	@Success	201		{object}	TaskDataResponse
//	@Failure	422		{object}	httpx.ValidationErrorResponse
//	@Router		/tasks [post]
func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	var req StoreTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	task, err := h.uc.Store(req)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.Data(w, http.StatusCreated, toTaskResponse(task))
}

// Update partially updates a task; omitted fields are left unchanged.
//
//	@Summary	Update a task
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		task	path		int					true	"Task ID"
//	@Param		body	body		UpdateTaskRequest	true	"Fields to update"
//	@Success	200		{object}	TaskDataResponse
//	@Failure	404		{object}	httpx.ErrorResponse
//	@Failure	422		{object}	httpx.ValidationErrorResponse
//	@Router		/tasks/{task} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := taskIDFromPath(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	var req UpdateTaskRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}
	if err := req.Validate(); err != nil {
		httpx.Error(w, err)
		return
	}

	task, err := h.uc.Update(id, req)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.Data(w, http.StatusOK, toTaskResponse(task))
}

// Destroy deletes a task.
//
//	@Summary	Delete a task
//	@Tags		tasks
//	@Param		task	path	int	true	"Task ID"
//	@Success	204
//	@Failure	404	{object}	httpx.ErrorResponse
//	@Router		/tasks/{task} [delete]
func (h *Handler) Destroy(w http.ResponseWriter, r *http.Request) {
	id, err := taskIDFromPath(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	if err := h.uc.Destroy(id); err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.NoContent(w)
}

func taskIDFromPath(r *http.Request) (uint, error) {
	raw := chi.URLParam(r, "task")

	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, apperror.New("invalid_parameter", "task id must be numeric", http.StatusBadRequest)
	}

	return uint(id), nil
}
