package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go-todo-app/db"
	"go-todo-app/views"

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
	h := &TodoHandler{Queries: queries}

	mux := http.NewServeMux()
	mux.Handle("GET /{$}", AppHandlerFunc(h.List))
	mux.Handle("POST /todos", AppHandlerFunc(h.Create))
	mux.Handle("PATCH /todos/{id}/toggle", AppHandlerFunc(h.Toggle))
	mux.Handle("DELETE /todos/{id}", AppHandlerFunc(h.Delete))

	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// TodoHandler handles HTTP requests for todo operations.
type TodoHandler struct {
	Queries *db.Queries
}

// List handles GET / — renders the full page or an HTMX partial.
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) error {
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")

	// Validate filter
	if filter != "" && filter != "todo" && filter != "done" {
		filter = ""
	}
	// Validate sort
	if sort != "" && sort != "description" && sort != "due_date" {
		sort = ""
	}

	todos, err := h.Queries.ListTodos(r.Context(), db.ListTodosParams{
		ByStatus:  filter != "",
		Status:    filter,
		SortField: sort,
	})
	if err != nil {
		return fmt.Errorf("listing todos: %w", err)
	}

	// HTMX requests get just the list partial; normal requests get the full page.
	if r.Header.Get("HX-Request") == "true" {
		return views.TodoList(todos, filter, sort).Render(r.Context(), w)
	} else {
		return views.TodoPage(todos, filter, sort, nil).Render(r.Context(), w)
	}
}

// Create handles POST /todos — adds a new todo and returns the updated list.
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return &HTTPError{Code: http.StatusBadRequest, Message: "Bad Request"}
	}

	description := strings.TrimSpace(r.FormValue("description"))
	dueDateStr := strings.TrimSpace(r.FormValue("due_date"))

	// Server-side validation
	var errors []string
	if description == "" {
		errors = append(errors, "Description is required.")
	}
	if len(description) > 500 {
		errors = append(errors, "Description must be 500 characters or fewer.")
	}
	if dueDateStr != "" {
		if _, err := time.Parse("2006-01-02", dueDateStr); err != nil {
			errors = append(errors, "Due date must be a valid date (YYYY-MM-DD).")
		}
	}

	if len(errors) > 0 {
		// Return error partial — swap into the form-errors div.
		w.Header().Set("HX-Retarget", "#form-errors")
		w.Header().Set("HX-Reswap", "innerHTML")
		return views.FormErrors(errors).Render(r.Context(), w)
	}

	// Build due_date as sql.NullString
	dueDate := sql.NullString{}
	if dueDateStr != "" {
		dueDate = sql.NullString{String: dueDateStr, Valid: true}
	}

	_, err := h.Queries.CreateTodo(r.Context(), db.CreateTodoParams{
		Description: description,
		Status:      "todo",
		DueDate:     dueDate,
	})
	if err != nil {
		return fmt.Errorf("creating todo: %w", err)
	}

	// Re-fetch and return the updated todo list.
	todos, err := h.Queries.ListTodos(r.Context(), db.ListTodosParams{})
	if err != nil {
		return fmt.Errorf("listing todos: %w", err)
	}

	// Clear any previous errors via OOB swap.
	w.Header().Set("HX-Trigger", "clearErrors")
	return views.TodoList(todos, "", "").Render(r.Context(), w)
}

// Delete handles DELETE /todos/{id} — removes a todo.
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return &HTTPError{Code: http.StatusBadRequest, Message: "Invalid todo ID"}
	}

	if err := h.Queries.DeleteTodo(r.Context(), id); err != nil {
		return fmt.Errorf("deleting todo %d: %w", id, err)
	}

	// Return empty body — HTMX outerHTML swap removes the element.
	w.WriteHeader(http.StatusOK)
	return nil
}

// Toggle handles PATCH /todos/{id}/toggle — toggles a todo's status.
func (h *TodoHandler) Toggle(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return &HTTPError{Code: http.StatusBadRequest, Message: "Invalid todo ID"}
	}

	todo, err := h.Queries.ToggleTodoStatus(r.Context(), id)
	if err != nil {
		return fmt.Errorf("toggling todo %d: %w", id, err)
	}

	// Return the updated todo item partial.
	return views.TodoItem(todo).Render(r.Context(), w)
}
