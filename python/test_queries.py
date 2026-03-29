import sqlite3
import tempfile
import unittest
from pathlib import Path

import aiosqlite

import queries


class QueriesTests(unittest.IsolatedAsyncioTestCase):
    async def asyncSetUp(self) -> None:
        self._temp_dir = tempfile.TemporaryDirectory()
        db_path = Path(self._temp_dir.name) / "todos.db"

        self.db = await aiosqlite.connect(db_path)
        self.db.row_factory = sqlite3.Row

        schema = Path(__file__).with_name("schema.sql").read_text()
        await self.db.executescript(schema)

    async def asyncTearDown(self) -> None:
        await self.db.close()
        self._temp_dir.cleanup()

    async def test_create_todo_sets_status_and_due_date(self) -> None:
        todo = await queries.create_todo(self.db, "write tests", "todo", "2026-04-01")

        self.assertEqual(todo.status, "todo")
        self.assertEqual(todo.due_date, "2026-04-01")

    async def test_list_todos_filters_done_and_sorts_description(self) -> None:
        await queries.create_todo(self.db, "zebra", "done", None)
        await queries.create_todo(self.db, "alpha", "done", None)
        await queries.create_todo(self.db, "middle", "todo", None)

        todos = await queries.list_todos(self.db, True, "done", "description")

        self.assertEqual(len(todos), 2)
        self.assertEqual(todos[0].description, "alpha")
        self.assertEqual(todos[1].description, "zebra")


if __name__ == "__main__":
    unittest.main()
