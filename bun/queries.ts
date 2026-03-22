import { Database } from "bun:sqlite";

export interface Todo {
  id: number;
  description: string;
  status: "todo" | "done";
  due_date: string | null;
  created_at: string;
}

export class Queries {
  constructor(private db: Database) {}

  listTodos(byStatus: boolean, status: string, sortField: string): Todo[] {
    return this.db
      .query<Todo, [boolean, string, string]>(
        `SELECT * FROM todos
         WHERE
           CASE WHEN ?1 THEN status = ?2 ELSE true END
           AND CASE WHEN ?3 != '' THEN true ELSE true END
         ORDER BY
           CASE ?3
             WHEN 'description' THEN description
             WHEN 'due_date' THEN due_date
             ELSE created_at
           END ASC`
      )
      .all(byStatus, status, sortField);
  }

  createTodo(description: string, status: string, dueDate: string | null): Todo {
    return this.db
      .query<Todo, [string, string, string | null]>(
        `INSERT INTO todos (description, status, due_date)
         VALUES (?1, ?2, ?3)
         RETURNING *`
      )
      .get(description, status, dueDate)!;
  }

  deleteTodo(id: number): void {
    this.db.query("DELETE FROM todos WHERE id = ?1").run(id);
  }

  toggleTodoStatus(id: number): Todo {
    return this.db
      .query<Todo, [number]>(
        `UPDATE todos
         SET status = CASE WHEN status = 'done' THEN 'todo' ELSE 'done' END
         WHERE id = ?1
         RETURNING *`
      )
      .get(id)!;
  }
}
