package main

import (
	"cototal/simple-crud/repos"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("CONNECTION_STRING")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /", func(wtr http.ResponseWriter, req *http.Request) {
		tasks, err := repos.GetAllTasks(db)
		if err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		json.NewEncoder(wtr).Encode(tasks)
	})

	router.HandleFunc("POST /", func(wtr http.ResponseWriter, req *http.Request) {
		var task repos.Task
		if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
			http.Error(wtr, err.Error(), http.StatusBadRequest)
			return
		}

		err = repos.CreateTask(db, &task)
		if err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
			return
		}
		wtr.Header().Set("Content-Type", "application/json")
		json.NewEncoder(wtr).Encode(task)
	})

	router.HandleFunc("GET /{id}", func(wtr http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			http.Error(wtr, "Invalid task ID", http.StatusBadRequest)
			return
		}

		task, err := repos.GetTask(db, id)

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
	})

	router.HandleFunc("PUT /{id}", func(wtr http.ResponseWriter, req *http.Request) {
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

		err = repos.UpdateTask(db, id, &task)
		if err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		json.NewEncoder(wtr).Encode(task)
	})

	router.HandleFunc("DELETE /{id}", func(wtr http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			http.Error(wtr, "Invalid task ID", http.StatusBadRequest)
			return
		}

		err = repos.DeleteTask(db, id)
		if err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
			return
		}
		wtr.WriteHeader(http.StatusNoContent)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting server on port :8080")
	server.ListenAndServe()
}
