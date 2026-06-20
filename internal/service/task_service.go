package service

import "github.com/abozorov/task_manager/internal/repo"

type TaskService struct {
	repo *repo.TaskRepo
}

func NewTaskService(repo *repo.TaskRepo) *TaskService {
	return &TaskService{
		repo: repo,
	}
}