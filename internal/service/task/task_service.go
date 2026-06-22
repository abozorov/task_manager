package service

import (
	"context"
	"fmt"

	"github.com/abozorov/task_manager/internal/models"
	taskRepo "github.com/abozorov/task_manager/internal/repo/task"
	userRepo "github.com/abozorov/task_manager/internal/repo/user"
	"github.com/abozorov/task_manager/pkg/errs"
)

type TaskService struct {
	userR *userRepo.UserRepo
	taskR *taskRepo.TaskRepo
}

func NewTaskService(userR *userRepo.UserRepo, taskR *taskRepo.TaskRepo) *TaskService {
	return &TaskService{
		userR: userR,
		taskR: taskR,
	}
}

func (t *TaskService) Create(ctx context.Context, task models.Task) error {
	// validation
	if !task.Validate(true) {
		return errs.ErrBadRequestBody
	}
	user, err := t.userR.GetByID(ctx, task.UserID)
	if err != nil {
		return fmt.Errorf("task_service.Create: %w", err)
	}

	// user is active
	if !user.DeletedAt.IsZero() {
		return fmt.Errorf("task_service.GetByID: %w", errs.ErrUserNotFound)
	}

	// creating
	err = t.taskR.Create(ctx, task)
	if err != nil {
		return fmt.Errorf("task_service.Create: %w", err)
	}

	return nil
}

func (t *TaskService) GetAll(ctx context.Context) ([]models.Task, error) {
	// get all Tasks
	allTasks, err := t.taskR.GetAll(ctx)
	if err != nil {
		return []models.Task{}, fmt.Errorf("task_service.GetAll: %w", err)
	}

	// get active Tasks
	activeTasks := make([]models.Task, 0, len(allTasks))
	for _, v := range allTasks {
		if v.DeletedAt.IsZero() {
			activeTasks = append(activeTasks, v)
		}
	}

	return activeTasks, nil
}

func (t *TaskService) GetByID(ctx context.Context, id int) (*models.Task, error) {
	// get all Tasks
	Task, err := t.taskR.GetByID(ctx, id)
	if err != nil {
		return &models.Task{}, fmt.Errorf("task_service.GetByID: %w", err)
	}

	// get active Tasks
	if !Task.DeletedAt.IsZero() {
		return &models.Task{}, fmt.Errorf("task_service.GetByID: %w", errs.ErrNotFound)
	}
	return Task, nil
}

func (t *TaskService) Update(ctx context.Context, Task models.Task) error {
	// validation
	if !Task.Validate(false) {
		return errs.ErrBadRequestBody
	}

	// updating
	err := t.taskR.Update(ctx, Task)
	if err != nil {
		return fmt.Errorf("Task_service.Update: %w", err)
	}
	return nil
}

func (t *TaskService) DeleteTask(ctx context.Context, id int) error {
	// delete Task
	err := t.taskR.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("Task_service.DeleteByID: %w", err)
	}
	return nil
}
