package taskHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abozorov/task_manager/internal/models"
	service "github.com/abozorov/task_manager/internal/service/task"
	"github.com/abozorov/task_manager/pkg/errs"
	"github.com/abozorov/task_manager/pkg/logger"
	"go.uber.org/zap"
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

type taskRequest struct {
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type taskResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

func newtaskResponse(t models.Task) *taskResponse {
	return &taskResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt.Format(time.RFC822Z),
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get user
	tsk := taskRequest{}
	err := json.NewDecoder(r.Body).Decode(&tsk)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.task -> task
	err = h.service.Create(r.Context(), *models.NewTask(
		tsk.UserID,
		tsk.Title,
		tsk.Description,
		tsk.Status,
	))
	if err != nil {
		h.logger.Error("task_handler.Create: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("task Created"))
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// load all
	tasks, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("task_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.User -> user
	resp := make([]taskResponse, 0, len(tasks))
	for _, v := range tasks {
		resp = append(resp, *newtaskResponse(v))
	}

	// write request
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("task_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("task_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// load all
	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("task_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.User -> user
	resp := *newtaskResponse(*task)

	// write request
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("task_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get user
	tsk := taskRequest{}
	err := json.NewDecoder(r.Body).Decode(&tsk)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.User -> user
	err = h.service.Update(r.Context(), models.Task{
		Title:       tsk.Title,
		Description: tsk.Description,
		Status:      tsk.Status,
	})
	if err != nil {
		h.logger.Error("task_handler.Update: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("task updated"))
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("task_handler.Delete: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// get by id
	err = h.service.DeleteTask(r.Context(), id)
	if err != nil {
		h.logger.Error("task_handler.Delete: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// write response
	w.Write([]byte(fmt.Sprintf("user %d deleted", id)))
}
