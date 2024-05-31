package main

import (
	"cototal/simple-crud/queries"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting server on port :8080")
	server.ListenAndServe()
}
