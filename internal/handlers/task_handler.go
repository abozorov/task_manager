package handlers

import (
	"net/http"

	"github.com/abozorov/task_manager/internal/service"
	"github.com/abozorov/task_manager/pkg/logger"
)

type TaskHandler struct {
	service *service.TaskService
	logger  *logger.Logger
}

func NewTaskHandler(service *service.TaskService, logger *logger.Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {}

