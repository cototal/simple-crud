package handlers

import (
	"database/sql"
	"net/http"
)

type MainRouter struct {
	db     *sql.DB
	router *http.ServeMux
}

func NewMainRouter(db *sql.DB, router *http.ServeMux) MainRouter {
	return MainRouter{db: db, router: router}
}

func (mr MainRouter) RunRoutes() {
	taskHandler := NewTaskHandler(mr.db)
	routes := map[string]http.HandlerFunc{
		"GET /":        taskHandler.GetTasks,
		"POST /":       taskHandler.CreateTask,
		"GET /{id}":    taskHandler.GetTask,
		"PUT /{id}":    taskHandler.UpdateTask,
		"DELETE /{id}": taskHandler.DeleteTask,
	}
	for pattern, handler := range routes {
		mr.router.HandleFunc(pattern, handler)
	}
}
