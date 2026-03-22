package handlers

import (
	"database/sql"
	"fmt"
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
		views.TodoList(todos, filter, sort).Render(r.Context(), w)
	} else {
		views.TodoPage(todos, filter, sort, nil).Render(r.Context(), w)
	}
	return nil
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
		views.FormErrors(errors).Render(r.Context(), w)
		return nil
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
	views.TodoList(todos, "", "").Render(r.Context(), w)
	return nil
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
	views.TodoItem(todo).Render(r.Context(), w)
	return nil
}
