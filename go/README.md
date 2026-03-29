# Go Todo App

A basic todo list web app built with Go, HTMX, and Basecoat UI.

## Stack

- **HTTP**: `net/http` (stdlib router, Go 1.22+)
- **Templates**: [templ](https://templ.guide/) (type-safe HTML templates)
- **Database**: [sqlc](https://sqlc.dev/) (SQL codegen) + [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go SQLite)
- **Frontend**: [HTMX](https://htmx.org/) + [Basecoat UI](https://basecoatui.com/) (via CDN)

## Prerequisites

- Go 1.22+
- [templ](https://templ.guide/) CLI: `go install github.com/a-h/templ/cmd/templ@latest`
- [sqlc](https://sqlc.dev/) CLI: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

## Code Generation

After modifying `.templ` files or `.sql` query files, regenerate the Go code:

```bash
sqlc generate
templ generate
```

## Run

```bash
go run .
```

The server starts at http://localhost:8080. The SQLite database (`data/todos.db`) is auto-created on first run.

## Type Check

```bash
go vet ./...
```

## Test

```bash
go test ./...
```

## Project Structure

```
go/
├── main.go              # Entrypoint: DB setup, router, server
├── sqlc.yaml            # sqlc configuration
├── data/                # SQLite database (auto-created, gitignored)
├── db/
│   ├── schema.sql       # Table schema
│   ├── query.sql        # sqlc-annotated SQL queries
│   ├── db.go            # Generated: DBTX interface, Queries struct
│   ├── models.go        # Generated: Todo struct
│   └── query.sql.go     # Generated: ListTodos, CreateTodo, DeleteTodo, ToggleTodoStatus
├── handlers/
│   └── todo.go          # HTTP handlers
└── views/
    ├── layout.templ      # HTML layout (head, body, CDN links)
    ├── todo.templ        # Todo components (page, form, list, item)
    ├── layout_templ.go   # Generated from layout.templ
    └── todo_templ.go     # Generated from todo.templ
```
