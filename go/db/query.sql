-- name: ListTodos :many
SELECT *
FROM todos
WHERE
    CASE WHEN sqlc.arg(by_status) THEN status = sqlc.arg(status) ELSE true END
    AND CASE WHEN sqlc.arg(sort_field) != '' THEN true ELSE true END
ORDER BY
    CASE ?3
        WHEN 'description' THEN description
        WHEN 'due_date' THEN due_date
        ELSE created_at
    END ASC;

-- name: CreateTodo :one
INSERT INTO todos (description, status, due_date)
VALUES (sqlc.arg(description), sqlc.arg(status), sqlc.arg(due_date))
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = sqlc.arg(id);

-- name: ToggleTodoStatus :one
UPDATE todos
SET status = CASE WHEN status = 'done' THEN 'todo' ELSE 'done' END
WHERE id = sqlc.arg(id)
RETURNING *;
