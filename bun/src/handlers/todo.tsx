import type { Queries } from "../db/queries";
import { TodoPage, TodoList, TodoItem, FormErrors } from "../views/todo";

function html(body: JSX.Element, status = 200, headers: Record<string, string> = {}): Response {
  return new Response(body as string, {
    status,
    headers: { "Content-Type": "text/html; charset=utf-8", ...headers },
  });
}

export class TodoHandler {
  constructor(private queries: Queries) {}

  /** GET / — renders the full page or an HTMX partial. */
  list(req: Request): Response {
    const url = new URL(req.url);
    let filter = url.searchParams.get("filter") ?? "";
    let sort = url.searchParams.get("sort") ?? "";

    // Validate filter
    if (filter !== "" && filter !== "todo" && filter !== "done") {
      filter = "";
    }
    // Validate sort
    if (sort !== "" && sort !== "description" && sort !== "due_date") {
      sort = "";
    }

    const todos = this.queries.listTodos({
      byStatus: filter !== "",
      status: filter,
      sortField: sort,
    });

    // HTMX requests get just the list partial; normal requests get the full page.
    if (req.headers.get("HX-Request") === "true") {
      return html(TodoList({ todos, filter, sort }));
    }
    return html(TodoPage({ todos, filter, sort, formErrors: [] }));
  }

  /** POST /todos — adds a new todo and returns the updated list. */
  async create(req: Request): Promise<Response> {
    const formData = await req.formData();
    const description = (formData.get("description") as string | null)?.trim() ?? "";
    const dueDateStr = (formData.get("due_date") as string | null)?.trim() ?? "";

    // Server-side validation
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
      // Return error partial — swap into the form-errors div.
      return html(FormErrors({ errors }), 200, {
        "HX-Retarget": "#form-errors",
        "HX-Reswap": "innerHTML",
      });
    }

    this.queries.createTodo({
      description,
      status: "todo",
      dueDate: dueDateStr || null,
    });

    // Re-fetch and return the updated todo list.
    const todos = this.queries.listTodos({ byStatus: false, status: "", sortField: "" });

    // Clear any previous errors via OOB swap.
    return html(TodoList({ todos, filter: "", sort: "" }), 200, {
      "HX-Trigger": "clearErrors",
    });
  }

  /** DELETE /todos/:id — removes a todo. */
  delete(req: Request, id: number): Response {
    this.queries.deleteTodo(id);
    // Return empty body — HTMX outerHTML swap removes the element.
    return new Response(null, { status: 200 });
  }

  /** PATCH /todos/:id/toggle — toggles a todo's status. */
  toggle(req: Request, id: number): Response {
    const todo = this.queries.toggleTodoStatus(id);
    return html(TodoItem({ todo }));
  }
}
