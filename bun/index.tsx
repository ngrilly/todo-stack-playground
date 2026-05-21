import { Database } from "bun:sqlite";
import { mkdirSync } from "node:fs";
import { Queries } from "./queries";
import { TodoPage, TodoList, TodoItem, FormErrors } from "./views";

// Ensure data directory exists.
mkdirSync("data", { recursive: true });

// Open SQLite database with WAL mode and foreign keys enabled.
const db = new Database("data/todos.db");
db.run("PRAGMA journal_mode = WAL");
db.run("PRAGMA foreign_keys = ON");

// Auto-create the schema on startup.
const ddl = await Bun.file(new URL("schema.sql", import.meta.url)).text();
db.run(ddl);

const queries = new Queries(db);

function html(
  body: JSX.Element,
  status = 200,
  headers: Record<string, string> = {},
): Response {
  return new Response(body as string, {
    status,
    headers: { "Content-Type": "text/html; charset=utf-8", ...headers },
  });
}

/** GET / — renders the full page or an HTMX partial. */
function list(req: Request): Response {
  const url = new URL(req.url);

  let filter = url.searchParams.get("filter") ?? "";
  if (filter !== "" && filter !== "todo" && filter !== "done") {
    filter = "";
  }

  let sort = url.searchParams.get("sort") ?? "";
  if (sort !== "" && sort !== "description" && sort !== "due_date") {
    sort = "";
  }

  const todos = queries.listTodos(filter !== "", filter, sort);

  if (req.headers.get("HX-Request") === "true") {
    return html(TodoList({ todos, filter, sort }));
  }
  return html(TodoPage({ todos, filter, sort, formErrors: [] }));
}

/** POST /todos — adds a new todo and returns the updated list. */
async function create(req: Request): Promise<Response> {
  const formData = await req.formData();
  const description =
    (formData.get("description") as string | null)?.trim() ?? "";
  const dueDateStr = (formData.get("due_date") as string | null)?.trim() ?? "";

  const errors: string[] = [];
  if (description === "") {
    errors.push("Description is required.");
  }
  if (description.length > 500) {
    errors.push("Description must be 500 characters or fewer.");
  }
  if (dueDateStr !== "") {
    const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
    if (!dateRegex.test(dueDateStr) || isNaN(new Date(dueDateStr).getTime())) {
      errors.push("Due date must be a valid date (YYYY-MM-DD).");
    }
  }

  if (errors.length > 0) {
    return html(FormErrors({ errors }), 200, {
      "HX-Retarget": "#form-errors",
      "HX-Reswap": "innerHTML",
    });
  }

  queries.createTodo(description, "todo", dueDateStr || null);

  const todos = queries.listTodos(false, "", "");
  return html(TodoList({ todos, filter: "", sort: "" }), 200, {
    "HX-Trigger": "clearErrors",
  });
}

/** DELETE /todos/:id — removes a todo. */
function remove(req: Request, id: number): Response {
  queries.deleteTodo(id);
  return new Response(null, { status: 200 });
}

/** PATCH /todos/:id/toggle — toggles a todo's status. */
function toggle(req: Request, id: number): Response {
  const todo = queries.toggleTodoStatus(id);
  return html(TodoItem({ todo }));
}

const server = Bun.serve({
  port: 8080,
  routes: {
    "/": {
      GET: (req) => list(req),
    },
    "/todos": {
      POST: (req) => create(req),
    },
    "/todos/:id/toggle": {
      PATCH: (req) => toggle(req, parseInt(req.params.id)),
    },
    "/todos/:id": {
      DELETE: (req) => remove(req, parseInt(req.params.id)),
    },
  },
  fetch() {
    return new Response("Not Found", { status: 404 });
  },
  error(error) {
    console.error("server error:", error);
    return new Response("Internal Server Error", { status: 500 });
  },
});

console.log(`Server listening on http://localhost:${server.port}`);
