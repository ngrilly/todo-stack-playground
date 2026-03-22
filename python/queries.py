from dataclasses import dataclass
from typing import Literal

import aiosqlite


@dataclass
class Todo:
    id: int
    description: str
    status: Literal["todo", "done"]
    due_date: str | None
    created_at: str


async def list_todos(
    db: aiosqlite.Connection,
    by_status: bool,
    status: str,
    sort_field: str,
) -> list[Todo]:
    async with db.execute(
        """
        SELECT * FROM todos
        WHERE
            CASE WHEN ? THEN status = ? ELSE true END
            AND CASE WHEN ? != '' THEN true ELSE true END
        ORDER BY
            CASE ?3
                WHEN 'description' THEN description
                WHEN 'due_date' THEN due_date
                ELSE created_at
            END ASC
        """,
        (by_status, status, sort_field),
    ) as cursor:
        return [Todo(**row) for row in await cursor.fetchall()]


async def create_todo(
    db: aiosqlite.Connection,
    description: str,
    status: str,
    due_date: str | None,
) -> Todo:
    async with db.execute(
        """
        INSERT INTO todos (description, status, due_date)
        VALUES (?, ?, ?)
        RETURNING *
        """,
        (description, status, due_date),
    ) as cursor:
        row = await cursor.fetchone()
        assert row is not None
        await db.commit()
        return Todo(**row)


async def delete_todo(db: aiosqlite.Connection, todo_id: int) -> None:
    await db.execute("DELETE FROM todos WHERE id = ?", (todo_id,))
    await db.commit()


async def toggle_todo_status(db: aiosqlite.Connection, todo_id: int) -> Todo:
    async with db.execute(
        """
        UPDATE todos
        SET status = CASE WHEN status = 'done' THEN 'todo' ELSE 'done' END
        WHERE id = ?
        RETURNING *
        """,
        (todo_id,),
    ) as cursor:
        row = await cursor.fetchone()
        assert row is not None
        await db.commit()
        return Todo(**row)
