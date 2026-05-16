using Microsoft.Data.Sqlite;

namespace TodoApp;

public class Queries(string connectionString)
{
    public async Task<IReadOnlyList<Todo>> ListTodosAsync(bool byStatus, string status, string sortField)
    {
        const string sql = """
            SELECT id, description, status, due_date, created_at
            FROM todos
            WHERE
              CASE WHEN $by_status THEN status = $status ELSE true END
              AND CASE WHEN $sort_field != '' THEN true ELSE true END
            ORDER BY
              CASE $sort_field
                WHEN 'description' THEN description
                WHEN 'due_date' THEN due_date
                ELSE created_at
              END ASC
            """;

        await using var connection = await OpenConnectionAsync();
        await using var command = connection.CreateCommand();
        command.CommandText = sql;
        command.Parameters.AddWithValue("$by_status", byStatus ? 1 : 0);
        command.Parameters.AddWithValue("$status", status);
        command.Parameters.AddWithValue("$sort_field", sortField);

        var todos = new List<Todo>();
        await using var reader = await command.ExecuteReaderAsync();
        while (await reader.ReadAsync())
        {
            todos.Add(ReadTodo(reader));
        }

        return todos;
    }

    public async Task<Todo> CreateTodoAsync(string description, string status, string? dueDate)
    {
        const string sql = """
            INSERT INTO todos (description, status, due_date)
            VALUES ($description, $status, $due_date)
            RETURNING id, description, status, due_date, created_at
            """;

        await using var connection = await OpenConnectionAsync();
        await using var command = connection.CreateCommand();
        command.CommandText = sql;
        command.Parameters.AddWithValue("$description", description);
        command.Parameters.AddWithValue("$status", status);
        command.Parameters.AddWithValue("$due_date", (object?)dueDate ?? DBNull.Value);

        await using var reader = await command.ExecuteReaderAsync();
        if (await reader.ReadAsync())
        {
            return ReadTodo(reader);
        }

        throw new InvalidOperationException("Failed to create todo.");
    }

    public async Task<Todo?> ToggleTodoStatusAsync(long id)
    {
        const string sql = """
            UPDATE todos
            SET status = CASE WHEN status = 'done' THEN 'todo' ELSE 'done' END
            WHERE id = $id
            RETURNING id, description, status, due_date, created_at
            """;

        await using var connection = await OpenConnectionAsync();
        await using var command = connection.CreateCommand();
        command.CommandText = sql;
        command.Parameters.AddWithValue("$id", id);

        await using var reader = await command.ExecuteReaderAsync();
        if (await reader.ReadAsync())
        {
            return ReadTodo(reader);
        }

        return null;
    }

    public async Task DeleteTodoAsync(long id)
    {
        const string sql = "DELETE FROM todos WHERE id = $id";

        await using var connection = await OpenConnectionAsync();
        await using var command = connection.CreateCommand();
        command.CommandText = sql;
        command.Parameters.AddWithValue("$id", id);
        await command.ExecuteNonQueryAsync();
    }

    private async Task<SqliteConnection> OpenConnectionAsync()
    {
        var connection = new SqliteConnection(connectionString);
        await connection.OpenAsync();

        await using var pragma = connection.CreateCommand();
        pragma.CommandText = "PRAGMA foreign_keys = ON";
        await pragma.ExecuteNonQueryAsync();

        return connection;
    }

    private static Todo ReadTodo(SqliteDataReader reader)
    {
        return new Todo(
            Id: reader.GetInt64(0),
            Description: reader.GetString(1),
            Status: reader.GetString(2),
            DueDate: reader.IsDBNull(3) ? null : reader.GetString(3),
            CreatedAt: reader.GetString(4)
        );
    }
}
