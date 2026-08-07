package task

import (
	"time"

	"github.com/st-man-hori/go-feature-template/internal/models"
)

// TaskResponse mirrors Resources/TaskResource.php: the single place that
// decides the API's on-the-wire shape for a Task, decoupled from the GORM
// model's column names/types (e.g. a nil DueDate becomes a JSON null, not
// an omitted key, matching Laravel's `$this->resource->due_date?->toDateString()`).
type TaskResponse struct {
	ID          uint    `json:"id" example:"1"`
	Title       string  `json:"title" example:"Write the quarterly report"`
	Description *string `json:"description" example:"Summarize Q2 sales figures"`
	DueDate     *string `json:"dueDate" example:"2026-08-01"`
	IsDone      bool    `json:"isDone" example:"false"`
	CreatedAt   string  `json:"createdAt" example:"2026-07-17T00:00:00Z"`
	UpdatedAt   string  `json:"updatedAt" example:"2026-07-17T00:00:00Z"`
}

// TaskDataResponse and TaskListDataResponse document the {"data": ...}
// envelope httpx.Data() builds at runtime. They exist only so swaggo's
// @Success annotations in handler.go have a concrete type to point at —
// handlers never construct these directly.
type TaskDataResponse struct {
	Data TaskResponse `json:"data"`
}

type TaskListDataResponse struct {
	Data []TaskResponse `json:"data"`
}

func toTaskResponse(t models.Task) TaskResponse {
	var dueDate *string
	if t.DueDate != nil {
		s := t.DueDate.Format("2006-01-02")
		dueDate = &s
	}

	return TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     dueDate,
		IsDone:      t.IsDone,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
}

func toTaskResponses(tasks []models.Task) []TaskResponse {
	out := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		out[i] = toTaskResponse(t)
	}
	return out
}
