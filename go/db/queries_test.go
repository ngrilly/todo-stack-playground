package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestQueries(t *testing.T) *Queries {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "todos.db")
	sqlDB, err := sql.Open("sqlite", "file:"+dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	if _, err := sqlDB.ExecContext(context.Background(), string(schema)); err != nil {
		t.Fatalf("init schema: %v", err)
	}

	return New(sqlDB)
}

func TestCreateTodo(t *testing.T) {
	q := newTestQueries(t)

	todo, err := q.CreateTodo(context.Background(), CreateTodoParams{
		Description: "write tests",
		Status:      "todo",
		DueDate:     sql.NullString{String: "2026-04-01", Valid: true},
	})
	if err != nil {
		t.Fatalf("create todo: %v", err)
	}

	if todo.Status != "todo" {
		t.Fatalf("unexpected status: got %q want %q", todo.Status, "todo")
	}
	if !todo.DueDate.Valid || todo.DueDate.String != "2026-04-01" {
		t.Fatalf("unexpected due_date: got %+v", todo.DueDate)
	}
}

func TestListTodosFilterAndSort(t *testing.T) {
	q := newTestQueries(t)

	_, err := q.CreateTodo(context.Background(), CreateTodoParams{
		Description: "zebra",
		Status:      "done",
		DueDate:     sql.NullString{},
	})
	if err != nil {
		t.Fatalf("create done todo zebra: %v", err)
	}

	_, err = q.CreateTodo(context.Background(), CreateTodoParams{
		Description: "alpha",
		Status:      "done",
		DueDate:     sql.NullString{},
	})
	if err != nil {
		t.Fatalf("create done todo alpha: %v", err)
	}

	_, err = q.CreateTodo(context.Background(), CreateTodoParams{
		Description: "middle",
		Status:      "todo",
		DueDate:     sql.NullString{},
	})
	if err != nil {
		t.Fatalf("create todo todo: %v", err)
	}

	todos, err := q.ListTodos(context.Background(), ListTodosParams{
		ByStatus:  true,
		Status:    "done",
		SortField: "description",
	})
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}

	if len(todos) != 2 {
		t.Fatalf("unexpected count: got %d want %d", len(todos), 2)
	}
	if todos[0].Description != "alpha" || todos[1].Description != "zebra" {
		t.Fatalf("unexpected order: got %q then %q", todos[0].Description, todos[1].Description)
	}
}
