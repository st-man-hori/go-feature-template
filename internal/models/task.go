// Package models holds GORM model structs shared across features — the
// equivalent of Laravel's app/Models/. Model definitions and table schema
// don't belong to any single Feature, so like the Laravel template, Task
// lives here instead of inside internal/features/task/.
package models

import "time"

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description *string
	DueDate     *time.Time `gorm:"type:date"`
	IsDone      bool       `gorm:"not null;default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Task) TableName() string {
	return "tasks"
}
