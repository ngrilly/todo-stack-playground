# C# Todo App

A basic todo list web app built with ASP.NET Core Minimal APIs, Razor Slices, HTMX, and Basecoat UI.

## Stack

- **HTTP**: ASP.NET Core Minimal APIs (.NET 10)
- **Templates**: [Razor Slices](https://github.com/DamianEdwards/RazorSlices) (typed Razor templates for minimal APIs)
- **Database**: [Microsoft.Data.Sqlite](https://learn.microsoft.com/dotnet/standard/data/sqlite/) + raw SQL
- **Frontend**: [HTMX](https://htmx.org/) + [Basecoat UI](https://basecoatui.com/) (via CDN)

## Prerequisites

- .NET 10 SDK (install with `brew install --cask dotnet-sdk`)

## Run

```bash
dotnet run --project TodoApp
```

The server starts at http://localhost:8080. The SQLite database (`data/todos.db`) is auto-created on first run.

## Type Check / Build

```bash
dotnet build
```

## Test

```bash
dotnet test
```
