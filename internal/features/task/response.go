package task

import (
	"time"

	"github.com/st-man-hori/go-feature-template/internal/models"
)

// TaskResponse は Resources/TaskResource.php に相当する: Task の API 上での
// 形を決める唯一の場所であり、GORM モデルのカラム名/型からは切り離されている
// (例えば DueDate が nil のときは、キーが省略されるのではなく JSON の null に
// なる。Laravel の `$this->resource->due_date?->toDateString()` と同じ挙動)。
type TaskResponse struct {
	ID          uint    `json:"id" example:"1"`
	Title       string  `json:"title" example:"Write the quarterly report"`
	Description *string `json:"description" example:"Summarize Q2 sales figures"`
	DueDate     *string `json:"dueDate" example:"2026-08-01"`
	IsDone      bool    `json:"isDone" example:"false"`
	CreatedAt   string  `json:"createdAt" example:"2026-07-17T00:00:00Z"`
	UpdatedAt   string  `json:"updatedAt" example:"2026-07-17T00:00:00Z"`
}

// TaskDataResponse と TaskListDataResponse は、httpx.Data() が実行時に組み立てる
// {"data": ...} エンベロープを表現するためだけの型。存在意義は、
// handler.go の swaggo の @Success アノテーションが指す先として具象型が
// 必要なことだけで、Handler がこれらを直接構築することはない。
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
