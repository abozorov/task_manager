package models

import (
	"strings"
	"time"
)

type Task struct {
	ID          int       `gorm:"serial;primaryKey"`
	Title       string    `gorm:"size:255"`
	Description string    `gorm:"size:255"`
	Status      string    `gorm:"size:100"`
	CreatedAt   time.Time `gorm:"not null"`
}

func (t *Task) Validate() bool {
	t.Title = strings.TrimSpace(t.Title)
	t.Description = strings.TrimSpace(t.Description)
	t.Status = strings.TrimSpace(t.Status)

	return len([]rune(t.Title)) > 0
}