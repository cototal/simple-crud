package main

import (
	"cototal/simple-crud/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	mux := http.NewServeMux()
	router := handlers.NewMainRouter(db, mux)
	router.RunRoutes()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on port :8080")
	server.ListenAndServe()
}
