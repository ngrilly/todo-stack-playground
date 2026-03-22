export interface Todo {
  id: number;
  description: string;
  status: "todo" | "done";
  due_date: string | null;
  created_at: string;
}
