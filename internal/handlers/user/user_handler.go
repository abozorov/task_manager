package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/abozorov/task_manager/internal/models"
	service "github.com/abozorov/task_manager/internal/service/user"
	"github.com/abozorov/task_manager/pkg/errs"
	"github.com/abozorov/task_manager/pkg/logger"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *service.UserService
	logger  *logger.Logger
}

func NewUserHandler(service *service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func newUser(u models.User) *user {
	return &user{
		ID:   u.ID,
		Name: u.Name,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get user
	usr := user{}
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.User -> user
	err = h.service.Create(r.Context(), *models.NewUser(usr.Name))
	if err != nil {
		h.logger.Error("user_handler.Create: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("User Created"))
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// load all
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("user_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.User -> user
	resp := make([]user, 0, len(users))
	for _, v := range users {
		resp = append(resp, *newUser(v))
	}

	// write request
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("user_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("user_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// get by id
	usr, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("user_handler.GetByID: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// transform models.User -> user
	resp := *newUser(*usr)

	// write response
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("user_handler.GetAll: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
	}
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get user
	usr := user{}
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		errs.ErrsToHttp(w, errs.ErrBadRequestBody)
		return
	}

	// creating & transform models.User -> user
	err = h.service.Update(r.Context(), models.User{
		ID:   usr.ID,
		Name: usr.Name,
	})
	if err != nil {
		h.logger.Error("user_handler.Update: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}
	w.Write([]byte("User updated"))
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// check path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("user_handler.Delete: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, errs.ErrBadRequest)
		return
	}

	// get by id
	err = h.service.DeleteUser(r.Context(), id)
	if err != nil {
		h.logger.Error("user_handler.Delete: ", zap.String("error", err.Error()))
		errs.ErrsToHttp(w, err)
		return
	}

	// write response
	w.Write([]byte(fmt.Sprintf("user %d deleteed", id)))
}
