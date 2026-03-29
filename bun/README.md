# Bun Todo App

A basic todo list web app built with Bun, HTMX, and Basecoat UI.

## Stack

- **HTTP**: [Bun.serve](https://bun.sh/) (built-in HTTP server)
- **Templates**: JSX via [@kitajs/html](https://html.kitajs.org/) (type-safe HTML generation)
- **Database**: [bun:sqlite](https://bun.sh/) (built-in SQLite) + TypeScript interfaces
- **Frontend**: [HTMX](https://htmx.org/) + [Basecoat UI](https://basecoatui.com/) (via CDN)

## Prerequisites

- [Bun](https://bun.sh/) 1.0+

## Install

```bash
bun install
```

## Run

```bash
bun run start
```

The server starts at http://localhost:8080. The SQLite database (`data/todos.db`) is auto-created on first run.

## Type Check

```bash
bun run typecheck
```

## Test

```bash
bun test
```

## Project Structure

```
bun/
├── package.json             # Dependencies and scripts
├── tsconfig.json            # TypeScript + JSX configuration
├── src/
│   ├── index.tsx            # Entrypoint: DB setup, router, server
│   ├── kita.d.ts            # HTMX type extensions for JSX
│   ├── db/
│   │   ├── schema.sql       # Table schema
│   │   ├── types.ts         # Todo interface
│   │   └── queries.ts       # SQL query functions (listTodos, createTodo, deleteTodo, toggleTodoStatus)
│   ├── handlers/
│   │   └── todo.tsx         # HTTP handlers
│   └── views/
│       ├── layout.tsx       # HTML layout (head, body, CDN links)
│       └── todo.tsx         # Todo components (page, form, list, item, filter bar, errors)
└── data/                    # SQLite database (auto-created, gitignored)
```
