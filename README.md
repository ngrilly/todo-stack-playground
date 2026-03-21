# Comparison of Go, Bun, and Python

Develop a basic todo list web app to compare Go, Bun, and Python.

Each version of the app should live in a dedicated subfolder.

For each version of the app, make sure that you:
- Provide a concise README, containing instructions on how to install and run the app.
- Type check the app.
- Run and test the app.

## Goal

The goal is to compare the three stacks on the following criteria:
- Simplicity
- Verbosity
- Speed of development
- Maintainability and stability (long-term support)
- Robustness
- Type-safety
- Integrated tooling
- Easy testing
- Easy deployment
- Libraries:
  - Input validation (from URL and from form, client-side and server-side)
  - Database access (raw SQL)
  - HTMX support (HTML generation, partials, reusable components accepting children components)

## Stack

### Frontend (shared across all versions)
- [HTMX](https://htmx.org/) (partial swaps for add/delete/filter/sort)
- [Tailwind CSS v4](https://tailwindcss.com/) (utility-first CSS framework, loaded via Play CDN)
- [Basecoat UI](https://basecoatui.com/) (component library built on Tailwind CSS)
- Client-side validation via HTML5 attributes

### Go (`go/`)
- net/http (stdlib router, Go 1.22+)
- [templ](https://templ.guide/) (type-safe HTML templates)
- [sqlc](https://sqlc.dev/) (SQL codegen for type-safe DB access)
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go SQLite driver)

### Bun (`bun/`)
- Bun.serve (built-in HTTP server)
- JSX (type-safe HTML generation)
- bun:sqlite (built-in SQLite) + TypeScript interfaces
- TypeScript for type checking

### Python (`python/`)
- [uv](https://docs.astral.sh/uv/) (package management)
- [ruff](https://docs.astral.sh/ruff/) (linting/formatting) + [ty](https://docs.astral.sh/ty/) (type checking)
- [FastAPI](https://fastapi.tiangolo.com/) (HTTP framework)
- [htpy](https://htpy.dev/) (type-safe HTML generation in pure Python)
- [aiosqlite](https://github.com/omnilib/aiosqlite) (async SQLite) + raw SQL + dataclasses

## Features
- List the todos
- Add a todo: description and optional due date (status defaults to "todo")
- Toggle a todo's status (mark as done/undone by clicking the checkbox)
- Delete a todo
- Filter todos (all, todo, done)
- Sort by description or due date

## Architecture
- Keep things as simple as possible.
- Single SQLite database per app, auto-created on startup.
- Single table: todos (id, description, status, due_date, created_at).
- HTMX partial swaps for all mutations and filters (no full page reloads).
- Server-side validation with error messages returned as HTML partials.

## UI
- Centered single-column layout (max-width ~672px) with padding.
- Add form: description input, optional due date picker, and "Add" button on one row.
- Filter/sort bar: Basecoat button groups for filter ("All", "Todo", "Done") and sort ("Description", "Due Date").
  Use Basecoat single compound class names (e.g. `btn-sm` for active, `btn-sm-outline` for inactive — never combine `btn` with `btn-outline`).
- Each todo is a single compact row: clickable checkbox (toggles done/undone via HTMX), description, optional due date, and icon delete button.
- Tailwind CSS, Basecoat UI, and HTMX loaded via CDN. No build step needed for frontend assets.

## TODO
- OpenAPI schema generation (FastAPI has it built-in; evaluate options for Go and Bun)
- Database migrations (evaluate goose, atlas, or similar tools)
- Live reload in development

## Alternatives we could try later

- Frontend:
  - [Datastar](https://data-star.dev/)
- Templates with static typing in Python:
  - [Ludic](https://getludic.dev/)
  - [FastHTML](https://fastht.ml/)
  - [Dominate](https://github.com/Knio/dominate)
- Templates with static typing in Go:
  - [htmgo](https://htmgo.dev/)
- HTTP framework for Bun:
  - [Hono](https://hono.dev/) (upgrade from Bun.serve if needed)
- Type-safe SQL in TS:
  - [Kysely](https://kysely.dev/) (type-safe SQL query builder)
  - [Drizzle](https://orm.drizzle.team/)
