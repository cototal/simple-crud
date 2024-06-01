package main

import (
	"cototal/simple-crud/queries"
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

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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
		rows, err := db.Query(queries.SelectAllTasks())
		if err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		tasks := make([]Task, 0, 20)
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Name); err != nil {
				http.Error(wtr, err.Error(), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		wtr.Header().Set("Content-Type", "application/json")
		json.NewEncoder(wtr).Encode(tasks)
	})

	router.HandleFunc("POST /", func(wtr http.ResponseWriter, req *http.Request) {
		var task Task
		if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
			http.Error(wtr, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := db.Exec(queries.CreateOneTask(), task.Name)
		if err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
			return
		}

		task.ID = int(id)
		wtr.Header().Set("Content-Type", "application/json")
		json.NewEncoder(wtr).Encode(task)
	})

	router.HandleFunc("GET /{id}", func(wtr http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			http.Error(wtr, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var task Task
		err = db.QueryRow(queries.SelectOneTask(), id).Scan(
			&task.ID, &task.Name)

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

		var task Task
		if err = json.NewDecoder(req.Body).Decode(&task); err != nil {
			http.Error(wtr, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec(queries.UpdateOneTask(), task.Name, id)

		if err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
			return
		}

		task.ID = id
		wtr.Header().Set("Content-Type", "application/json")
		json.NewEncoder(wtr).Encode(task)
	})

	router.HandleFunc("DELETE /{id}", func(wtr http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			http.Error(wtr, "Invalid task ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(queries.DeleteOneTask(), id)
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
