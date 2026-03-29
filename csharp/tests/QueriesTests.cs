using System.Text;
using CSharpTodoApp;
using Xunit;

namespace CSharpTodoApp.Tests;

public class QueriesTests
{
    private static readonly string SchemaPath = System.IO.Path.GetFullPath(
        System.IO.Path.Combine(AppContext.BaseDirectory, "../../../../schema.sql")
    );

    [Fact]
    public async Task CreateTodoAsync_SetsStatusAndDueDate()
    {
        using var tempDir = new TempDirectory();
        var queries = await CreateQueriesAsync(tempDir.Path);

        var todo = await queries.CreateTodoAsync("write tests", "todo", "2026-04-01");

        Assert.Equal("todo", todo.Status);
        Assert.Equal("2026-04-01", todo.DueDate);
    }

    [Fact]
    public async Task ListTodosAsync_FiltersDoneAndSortsByDescription()
    {
        using var tempDir = new TempDirectory();
        var queries = await CreateQueriesAsync(tempDir.Path);

        await queries.CreateTodoAsync("zebra", "done", null);
        await queries.CreateTodoAsync("alpha", "done", null);
        await queries.CreateTodoAsync("middle", "todo", null);

        var todos = await queries.ListTodosAsync(byStatus: true, status: "done", sortField: "description");

        Assert.Equal(2, todos.Count);
        Assert.Equal("alpha", todos[0].Description);
        Assert.Equal("zebra", todos[1].Description);
    }

    private static async Task<Queries> CreateQueriesAsync(string tempDirPath)
    {
        Directory.CreateDirectory(tempDirPath);
        var dbPath = Path.Combine(tempDirPath, "todos.db");
        var connectionString = $"Data Source={dbPath}";

        await using var connection = new Microsoft.Data.Sqlite.SqliteConnection(connectionString);
        await connection.OpenAsync();

        var schemaSql = await File.ReadAllTextAsync(SchemaPath, Encoding.UTF8);
        await using var command = connection.CreateCommand();
        command.CommandText = schemaSql;
        await command.ExecuteNonQueryAsync();

        return new Queries(connectionString);
    }

    private class TempDirectory : IDisposable
    {
        public string Path { get; } = System.IO.Path.Combine(System.IO.Path.GetTempPath(), $"csharp-todo-test-{Guid.NewGuid():N}");

        public TempDirectory()
        {
            Directory.CreateDirectory(Path);
        }

        public void Dispose()
        {
            if (Directory.Exists(Path))
            {
                Directory.Delete(Path, recursive: true);
            }
        }
    }
}
