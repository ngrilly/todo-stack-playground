package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-todo-app/db"
	"go-todo-app/views"
)

// TodoHandler handles HTTP requests for todo operations.
type TodoHandler struct {
	Queries *db.Queries
}

// List handles GET / — renders the full page or an HTMX partial.
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("error listing todos: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// HTMX requests get just the list partial; normal requests get the full page.
	if r.Header.Get("HX-Request") == "true" {
		views.TodoList(todos, filter, sort).Render(r.Context(), w)
	} else {
		views.TodoPage(todos, filter, sort, nil).Render(r.Context(), w)
	}
}

// Create handles POST /todos — adds a new todo and returns the updated list.
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
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
		views.FormErrors(errors).Render(r.Context(), w)
		return
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
		log.Printf("error creating todo: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Re-fetch and return the updated todo list.
	todos, err := h.Queries.ListTodos(r.Context(), db.ListTodosParams{})
	if err != nil {
		log.Printf("error listing todos: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Clear any previous errors via OOB swap.
	w.Header().Set("HX-Trigger", "clearErrors")
	views.TodoList(todos, "", "").Render(r.Context(), w)
}

// Delete handles DELETE /todos/{id} — removes a todo.
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := h.Queries.DeleteTodo(r.Context(), id); err != nil {
		log.Printf("error deleting todo %d: %v", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return empty body — HTMX outerHTML swap removes the element.
	w.WriteHeader(http.StatusOK)
}

// Toggle handles PATCH /todos/{id}/toggle — toggles a todo's status.
func (h *TodoHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.Queries.ToggleTodoStatus(r.Context(), id)
	if err != nil {
		log.Printf("error toggling todo %d: %v", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return the updated todo item partial.
	views.TodoItem(todo).Render(r.Context(), w)
}
