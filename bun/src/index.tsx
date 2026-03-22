import { Database } from "bun:sqlite";
import { mkdirSync } from "node:fs";
import { Queries } from "./db/queries";
import { TodoHandler } from "./handlers/todo";

// Ensure data directory exists.
mkdirSync("data", { recursive: true });

// Open SQLite database with WAL mode and foreign keys enabled.
const db = new Database("data/todos.db");
db.run("PRAGMA journal_mode = WAL");
db.run("PRAGMA foreign_keys = ON");

// Auto-create the schema on startup.
const ddl = await Bun.file(new URL("db/schema.sql", import.meta.url)).text();
db.run(ddl);

const queries = new Queries(db);
const handler = new TodoHandler(queries);

const server = Bun.serve({
  port: 8080,
  routes: {
    "/": {
      GET: (req) => handler.list(req),
    },
    "/todos": {
      POST: (req) => handler.create(req),
    },
    "/todos/:id/toggle": {
      PATCH: (req) => handler.toggle(req, parseInt(req.params.id)),
    },
    "/todos/:id": {
      DELETE: (req) => handler.delete(req, parseInt(req.params.id)),
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
