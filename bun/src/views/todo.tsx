import type { Todo } from "../db/types";
import { Layout } from "./layout";

function filterBtnClass(active: boolean): string {
  return active ? "btn-sm" : "btn-sm-outline";
}

export function TodoPage({
  todos,
  filter,
  sort,
  formErrors,
}: {
  todos: Todo[];
  filter: string;
  sort: string;
  formErrors: string[];
}) {
  return (
    <Layout title="Todo App">
      <h1 class="text-2xl font-bold mb-6">Todo List</h1>
      <div id="form-errors">
        <FormErrors errors={formErrors} />
      </div>
      <TodoForm />
      <div id="todo-list">
        <TodoList todos={todos} filter={filter} sort={sort} />
      </div>
    </Layout>
  );
}

export function FormErrors({ errors }: { errors: string[] }) {
  if (errors.length === 0) return <></>;
  return (
    <div class="alert alert-destructive mb-4" role="alert">
      <ul class="list-disc list-inside">
        {errors.map((err) => (
          <li safe>{err}</li>
        ))}
      </ul>
    </div>
  );
}

export function TodoForm() {
  return (
    <form hx-post="/todos" hx-target="#todo-list" hx-swap="innerHTML" class="card p-4 mb-4">
      <div class="flex items-end gap-3">
        <div class="flex-1 flex flex-col gap-1.5">
          <label class="label" for="description">
            Description
          </label>
          <input
            class="input w-full"
            type="text"
            id="description"
            name="description"
            placeholder="What needs to be done?"
            required
            minlength={1}
            maxlength={500}
          />
        </div>
        <div class="flex flex-col gap-1.5">
          <label class="label" for="due_date">
            Due Date
          </label>
          <input class="input" type="date" id="due_date" name="due_date" />
        </div>
        <button class="btn" type="submit">
          Add
        </button>
      </div>
    </form>
  );
}

export function FilterBar({ currentFilter, currentSort }: { currentFilter: string; currentSort: string }) {
  return (
    <div class="flex flex-wrap items-center gap-4 mb-4">
      <span class="text-sm font-medium text-muted-foreground">Filter:</span>
      <div role="group" class="button-group">
        <button
          class={filterBtnClass(currentFilter === "")}
          hx-get={`/?filter=&sort=${currentSort}`}
          hx-target="#todo-list"
          hx-swap="innerHTML"
        >
          All
        </button>
        <button
          class={filterBtnClass(currentFilter === "todo")}
          hx-get={`/?filter=todo&sort=${currentSort}`}
          hx-target="#todo-list"
          hx-swap="innerHTML"
        >
          Todo
        </button>
        <button
          class={filterBtnClass(currentFilter === "done")}
          hx-get={`/?filter=done&sort=${currentSort}`}
          hx-target="#todo-list"
          hx-swap="innerHTML"
        >
          Done
        </button>
      </div>
      <span class="text-sm font-medium text-muted-foreground">Sort:</span>
      <div role="group" class="button-group">
        <button
          class={filterBtnClass(currentSort === "description")}
          hx-get={`/?filter=${currentFilter}&sort=description`}
          hx-target="#todo-list"
          hx-swap="innerHTML"
        >
          Description
        </button>
        <button
          class={filterBtnClass(currentSort === "due_date")}
          hx-get={`/?filter=${currentFilter}&sort=due_date`}
          hx-target="#todo-list"
          hx-swap="innerHTML"
        >
          Due Date
        </button>
      </div>
    </div>
  );
}

export function TodoList({ todos, filter, sort }: { todos: Todo[]; filter: string; sort: string }) {
  return (
    <>
      <FilterBar currentFilter={filter} currentSort={sort} />
      {todos.length === 0 ? (
        <p class="text-muted-foreground text-center py-8">No todos found.</p>
      ) : (
        <div class="flex flex-col gap-1">
          {todos.map((todo) => (
            <TodoItem todo={todo} />
          ))}
        </div>
      )}
    </>
  );
}

export function TodoItem({ todo }: { todo: Todo }) {
  return (
    <div class="flex items-center gap-3 px-2 py-1.5 rounded-md hover:bg-muted/50" id={`todo-${todo.id}`}>
      <input
        type="checkbox"
        class="checkbox"
        checked={todo.status === "done"}
        hx-patch={`/todos/${todo.id}/toggle`}
        hx-target={`#todo-${todo.id}`}
        hx-swap="outerHTML"
      />
      {todo.status === "done" ? (
        <span class="flex-1 min-w-0 truncate line-through text-muted-foreground" safe>
          {todo.description}
        </span>
      ) : (
        <span class="flex-1 min-w-0 truncate" safe>
          {todo.description}
        </span>
      )}
      {todo.due_date && (
        <span class="text-sm text-muted-foreground whitespace-nowrap" safe>
          {todo.due_date}
        </span>
      )}
      <button
        class="btn-sm-icon-destructive"
        hx-delete={`/todos/${todo.id}`}
        hx-target={`#todo-${todo.id}`}
        hx-swap="outerHTML"
        hx-confirm="Delete this todo?"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M3 6h18"></path>
          <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"></path>
          <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"></path>
        </svg>
      </button>
    </div>
  );
}
