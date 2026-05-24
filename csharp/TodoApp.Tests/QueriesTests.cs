using TodoApp;
using Xunit;

namespace TodoApp.Tests;

public class QueriesTests : IAsyncLifetime
{
    private static readonly string SchemaPath = Path.GetFullPath(
        Path.Combine(AppContext.BaseDirectory, "../../../../TodoApp/schema.sql")
    );

    private readonly string _tempDirPath;
    private Queries _queries = null!;

    public QueriesTests()
    {
        _tempDirPath = Path.Combine(Path.GetTempPath(), $"csharp-todo-test-{Guid.NewGuid():N}");
        Directory.CreateDirectory(_tempDirPath);
    }

    public async Task InitializeAsync()
    {
        var dbPath = Path.Combine(_tempDirPath, "todos.db");

        await using var connection = new Microsoft.Data.Sqlite.SqliteConnection($"Data Source={dbPath}");
        await connection.OpenAsync();

        var schemaSql = await File.ReadAllTextAsync(SchemaPath);
        await using var command = connection.CreateCommand();
        command.CommandText = schemaSql;
        await command.ExecuteNonQueryAsync();

        _queries = new Queries($"Data Source={dbPath}");
    }

    public Task DisposeAsync()
    {
        Directory.Delete(_tempDirPath, recursive: true);
        return Task.CompletedTask;
    }

    [Fact]
    public async Task CreateTodoAsync_SetsStatusAndDueDate()
    {
        var todo = await _queries.CreateTodoAsync("write tests", "todo", "2026-04-01");

        Assert.Equal("todo", todo.Status);
        Assert.Equal("2026-04-01", todo.DueDate);
    }

    [Fact]
    public async Task ListTodosAsync_FiltersDoneAndSortsByDescription()
    {
        await _queries.CreateTodoAsync("zebra", "done", null);
        await _queries.CreateTodoAsync("alpha", "done", null);
        await _queries.CreateTodoAsync("middle", "todo", null);

        var todos = await _queries.ListTodosAsync(byStatus: true, status: "done", sortField: "description");

        Assert.Equal(2, todos.Count);
        Assert.Equal("alpha", todos[0].Description);
        Assert.Equal("zebra", todos[1].Description);
    }
}
