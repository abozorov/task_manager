package repo

import (
	"context"

	"github.com/abozorov/task_manager/internal/models"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{
		db: db,
	}
}

func (t *TaskRepo) Create(ctx context.Context, task models.Task) error {
	tx := t.db.WithContext(ctx).Create(&task)
	return tx.Error
}

func (t *TaskRepo) GetAll(ctx context.Context) ([]models.Task, error) {
	tasks := make([]models.Task, 0)
	tx := t.db.WithContext(ctx).Select("id", "title", "description", "status", "created_at", "is_active").Find(&tasks)

	return tasks, tx.Error
}

func (t *TaskRepo) GetByID(ctx context.Context, id int) (*models.Task, error) {
	return &models.Task{}, nil
}
