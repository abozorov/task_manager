package repo

import (
	"context"
	"time"

	"github.com/abozorov/task_manager/internal/models"
	"github.com/abozorov/task_manager/pkg/errs"
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
	return errs.PostgresToErrs(tx.Error)
}

func (t *TaskRepo) GetAll(ctx context.Context) ([]models.Task, error) {
	tasks := make([]models.Task, 0)
	tx := t.db.WithContext(ctx).Select(
		"id",
		"user_id",
		"title",
		"description",
		"status",
		"created_at",
		"deleted_at",
	).Find(&tasks)

	return tasks, errs.PostgresToErrs(tx.Error)
}

func (t *TaskRepo) GetByID(ctx context.Context, id int) (*models.Task, error) {
	task := &models.Task{}
	tx := t.db.WithContext(ctx).Select(
		"id",
		"user_id",
		"title",
		"description",
		"status",
		"created_at",
		"deleted_at",
	).Where("id = ", id).First(&task)
	return task, errs.PostgresToErrs(tx.Error)
}

func (t *TaskRepo) Update(ctx context.Context, task models.Task) error {
	tx := t.db.WithContext(ctx).Save(&task)
	return errs.PostgresToErrs(tx.Error)
}

func (t *TaskRepo) DeleteByID(ctx context.Context, id int) error {
	task, err := t.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !task.DeletedAt.IsZero() {
		return errs.ErrNotFound
	}
	tx := t.db.WithContext(ctx).Model(&task).Update("deleted_at", time.Now())
	return errs.PostgresToErrs(tx.Error)
}

func (t *TaskRepo) DeleteByUserID(ctx context.Context, userID int) error {
	tx := t.db.WithContext(ctx).Table("tasks").Where(
		"user_id = ? AND deleted_at IS null",
		userID,
	).Update("deleted_at", time.Now())

	return errs.PostgresToErrs(tx.Error)
}
