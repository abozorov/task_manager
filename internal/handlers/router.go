package handlers

import (
	"net/http"

	"github.com/abozorov/task_manager/internal/handlers/middleware"
	taskHandler "github.com/abozorov/task_manager/internal/handlers/task"
	userHandler "github.com/abozorov/task_manager/internal/handlers/user"
)

type Router struct {
	*http.ServeMux
}

func NewRouter(t *taskHandler.TaskHandler, u *userHandler.UserHandler) *Router {
	mux := http.NewServeMux()

	// TASK
	mux.Handle("POST   /task", middleware.Logging(http.HandlerFunc(t.Create)))
	mux.Handle("GET    /tasks", middleware.Logging(http.HandlerFunc(t.GetAll)))
	mux.Handle("GET    /task/{id}", middleware.Logging(http.HandlerFunc(t.GetByID)))
	mux.Handle("PUT    /task", middleware.Logging(http.HandlerFunc(t.Update)))
	mux.Handle("DELETE /task/{id}", middleware.Logging(http.HandlerFunc(t.Delete)))

	// USER
	mux.Handle("POST /users", middleware.Logging(http.HandlerFunc(u.Create)))
	mux.Handle("GET /users", middleware.Logging(http.HandlerFunc(u.GetAll)))
	mux.Handle("GET /user/{id}", middleware.Logging(http.HandlerFunc(u.GetByID)))
	mux.Handle("PUT /user", middleware.Logging(http.HandlerFunc(u.Update)))
	mux.Handle("DELETE /user/{id}", middleware.Logging(http.HandlerFunc(u.Delete)))

	return &Router{
		mux,
	}

}

/*
POST   /tasks
GET    /tasks
GET    /tasks/{id}
PUT    /tasks/{id}
DELETE /tasks/{id}


GET /tasks?status=done
GET /tasks?page=1&limit=10
*/
