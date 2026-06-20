package handlers

import "net/http"

type Router struct {
	*http.ServeMux
}

func NewRouter(h TaskHandler) *Router {
	taskMux := http.NewServeMux()

	taskMux.HandleFunc("POST   /tasks", h.Create)
	taskMux.HandleFunc("GET    /tasks", h.GetAll)
	taskMux.HandleFunc("GET    /tasks/{id}", h.GetByID)
	taskMux.HandleFunc("PUT    /tasks/{id}", h.Update)
	taskMux.HandleFunc("DELETE /tasks/{id}", h.Delete)

	return &Router{
		taskMux,
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
