import logging
import re
import sqlite3
from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager
from pathlib import Path

import aiosqlite
from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse, Response

import queries
import views

logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None]:
    # Ensure data directory exists.
    Path("data").mkdir(exist_ok=True)

    # Open SQLite database with WAL mode and foreign keys enabled.
    db = await aiosqlite.connect("data/todos.db")
    db.row_factory = sqlite3.Row
    await db.execute("PRAGMA journal_mode = WAL")
    await db.execute("PRAGMA foreign_keys = ON")

    # Auto-create the schema on startup.
    ddl = Path(__file__).with_name("schema.sql").read_text()
    await db.executescript(ddl)

    app.state.db = db

    yield

    await db.close()


app = FastAPI(lifespan=lifespan)


def _get_db(request: Request) -> aiosqlite.Connection:
    return request.app.state.db


@app.get("/", response_class=HTMLResponse)
async def list_todos(request: Request) -> Response:
    db = _get_db(request)
    params = request.query_params

    filter = params.get("filter", "")
    if filter not in ("", "todo", "done"):
        filter = ""

    sort = params.get("sort", "")
    if sort not in ("", "description", "due_date"):
        sort = ""

    todos = await queries.list_todos(db, filter != "", filter, sort)

    # HTMX requests get just the list partial; normal requests get the full page.
    if request.headers.get("HX-Request") == "true":
        html = views.todo_list(todos, filter, sort)
    else:
        html = views.todo_page(todos, filter, sort, [])

    return HTMLResponse(str(html))


@app.post("/todos", response_class=HTMLResponse)
async def create_todo(request: Request) -> Response:
    db = _get_db(request)
    form_data = await request.form()
    description = str(form_data.get("description") or "").strip()
    due_date_str = str(form_data.get("due_date") or "").strip()

    # Server-side validation.
    errors: list[str] = []
    if not description:
        errors.append("Description is required.")
    if len(description) > 500:
        errors.append("Description must be 500 characters or fewer.")
    if due_date_str:
        if not re.fullmatch(r"\d{4}-\d{2}-\d{2}", due_date_str):
            errors.append("Due date must be a valid date (YYYY-MM-DD).")

    if errors:
        # Return error partial — swap into the form-errors div.
        return HTMLResponse(
            str(views.form_errors_view(errors)),
            headers={
                "HX-Retarget": "#form-errors",
                "HX-Reswap": "innerHTML",
            },
        )

    await queries.create_todo(db, description, "todo", due_date_str or None)

    # Re-fetch and return the updated todo list.
    todos = await queries.list_todos(db, False, "", "")

    return HTMLResponse(
        str(views.todo_list(todos, "", "")),
        headers={"HX-Trigger": "clearErrors"},
    )


@app.patch("/todos/{todo_id}/toggle", response_class=HTMLResponse)
async def toggle_todo(request: Request, todo_id: int) -> Response:
    db = _get_db(request)
    todo = await queries.toggle_todo_status(db, todo_id)
    return HTMLResponse(str(views.todo_item(todo)))


@app.delete("/todos/{todo_id}")
async def delete_todo(request: Request, todo_id: int) -> Response:
    db = _get_db(request)
    await queries.delete_todo(db, todo_id)
    return Response(status_code=200)
