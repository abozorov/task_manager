package models

import (
	"strings"
	"time"
)

type Task struct {
	ID          int       `gorm:"serial;primaryKey"`
	UserID      int       `gorm:"not null;check:user_id>0;"`
	Title       string    `gorm:"size:255"`
	Description string    `gorm:"size:255"`
	Status      string    `gorm:"size:100"`
	CreatedAt   time.Time `gorm:"default:current_timestamp;"`
	DeletedAt   time.Time `gorm:"default:null;"`
}

func NewTask(title, description, status string) *Task {
	return &Task{
		Title:       title,
		Description: description,
		Status:      status,
	}
}

func (t *Task) Validate(create bool) bool {
	t.Title = strings.TrimSpace(t.Title)
	t.Description = strings.TrimSpace(t.Description)
	t.Status = strings.TrimSpace(t.Status)

	return len([]rune(t.Title)) > 0 &&
		len([]rune(t.Description)) > 0 &&
		len([]rune(t.Status)) > 0 &&
		(create || t.ID > 0)
}
