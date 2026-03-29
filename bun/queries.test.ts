import { Database } from "bun:sqlite";
import { beforeEach, describe, expect, it } from "bun:test";
import { mkdtempSync, readFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { Queries } from "./queries";

const schema = readFileSync(new URL("./schema.sql", import.meta.url), "utf8");

let queries: Queries;

beforeEach(() => {
  const tempDir = mkdtempSync(join(tmpdir(), "bun-todo-test-"));
  const db = new Database(join(tempDir, "todos.db"));
  db.run(schema);
  queries = new Queries(db);
});

describe("Queries", () => {
  it("creates a todo with default todo status and due date", () => {
    const todo = queries.createTodo("write tests", "todo", "2026-04-01");

    expect(todo.status).toBe("todo");
    expect(todo.due_date).toBe("2026-04-01");
  });

  it("lists todos filtered by done and sorted by description", () => {
    queries.createTodo("zebra", "done", null);
    queries.createTodo("alpha", "done", null);
    queries.createTodo("middle", "todo", null);

    const todos = queries.listTodos(true, "done", "description");

    expect(todos).toHaveLength(2);
    expect(todos[0]?.description).toBe("alpha");
    expect(todos[1]?.description).toBe("zebra");
  });
});
