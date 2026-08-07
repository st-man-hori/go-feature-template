// Package models は、複数の Feature にまたがって共有される GORM モデルの構造体を
// 保持する — Laravel の app/Models/ に相当する。モデル定義とテーブルスキーマは
// 特定の Feature に属するものではないため、Laravel版と同じく Task は
// internal/features/task/ の中ではなくここに置く。
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
