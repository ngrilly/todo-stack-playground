package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"net/http"
	"os"

	"go-todo-app/db"
	"go-todo-app/handlers"

	_ "modernc.org/sqlite"
)

//go:embed db/schema.sql
var ddl string

func main() {
	// Ensure data directory exists.
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	// Open SQLite database with WAL mode and foreign keys enabled.
	sqlDB, err := sql.Open("sqlite", "file:data/todos.db?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer sqlDB.Close()

	// Auto-create the schema on startup.
	if _, err := sqlDB.ExecContext(context.Background(), ddl); err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}

	queries := db.New(sqlDB)
	h := &handlers.TodoHandler{Queries: queries}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", h.List)
	mux.HandleFunc("POST /todos", h.Create)
	mux.HandleFunc("PATCH /todos/{id}/toggle", h.Toggle)
	mux.HandleFunc("DELETE /todos/{id}", h.Delete)

	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
