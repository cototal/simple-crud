package handlers

import (
	"cototal/simple-crud/repos"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	db *sql.DB
}

func NewTaskHandler(db *sql.DB) TaskHandler {
	return TaskHandler{db: db}
}

func (handler TaskHandler) GetTasks(wtr http.ResponseWriter, req *http.Request) {
	tasks, err := repos.GetAllTasks(handler.db)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wtr).Encode(tasks)
}

func (handler TaskHandler) CreateTask(wtr http.ResponseWriter, req *http.Request) {
	var task repos.Task
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		http.Error(wtr, err.Error(), http.StatusBadRequest)
		return
	}

	err := repos.CreateTask(handler.db, &task)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}
	wtr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wtr).Encode(task)
}

func (handler TaskHandler) GetTask(wtr http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(wtr, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := repos.GetTask(handler.db, id)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(wtr, "Task not found", http.StatusNotFound)
		} else {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wtr).Encode(task)
}

func (handler TaskHandler) UpdateTask(wtr http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(wtr, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task repos.Task
	if err = json.NewDecoder(req.Body).Decode(&task); err != nil {
		http.Error(wtr, err.Error(), http.StatusBadRequest)
		return
	}

	err = repos.UpdateTask(handler.db, id, &task)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wtr).Encode(task)
}

func (handler TaskHandler) DeleteTask(wtr http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(wtr, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = repos.DeleteTask(handler.db, id)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}
	wtr.WriteHeader(http.StatusNoContent)
}
