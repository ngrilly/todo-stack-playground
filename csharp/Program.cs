using System.Globalization;
using System.Text.RegularExpressions;
using CSharpTodoApp;

var builder = WebApplication.CreateBuilder(args);
builder.WebHost.UseUrls("http://localhost:8080");

var dataDir = Path.Combine(builder.Environment.ContentRootPath, "data");
var dbPath = Path.Combine(dataDir, "todos.db");
var connectionString = $"Data Source={dbPath}";

builder.Services.AddSingleton(new Queries(connectionString));

var app = builder.Build();

InitializeDatabase(dataDir, dbPath, app.Environment.ContentRootPath);

app.Use(async (context, next) =>
{
    try
    {
        await next();
    }
    catch (Exception ex)
    {
        app.Logger.LogError(ex, "Unhandled server error");

        if (!context.Response.HasStarted)
        {
            context.Response.StatusCode = StatusCodes.Status500InternalServerError;
            context.Response.ContentType = "text/plain; charset=utf-8";
            await context.Response.WriteAsync("Internal Server Error");
        }
    }
});

app.MapGet("/", (Func<HttpContext, Task<IResult>>)(async context =>
{
    var request = context.Request;
    var queries = context.RequestServices.GetRequiredService<Queries>();
    var filter = request.Query["filter"].ToString();
    var sort = request.Query["sort"].ToString();

    if (filter is not ("" or "todo" or "done"))
    {
        filter = "";
    }

    if (sort is not ("" or "description" or "due_date"))
    {
        sort = "";
    }

    var todos = await queries.ListTodosAsync(filter != "", filter, sort);

    if (IsHtmxRequest(request))
    {
        return Results.Extensions.RazorSlice<CSharpTodoApp.Slices.TodoList, TodoListModel>(
            new TodoListModel(todos, filter, sort)
        );
    }

    return Results.Extensions.RazorSlice<CSharpTodoApp.Slices.TodoPage, TodoPageModel>(
        new TodoPageModel(todos, filter, sort, [])
    );
}));

app.MapPost("/todos", (Func<HttpContext, Task<IResult>>)(async context =>
{
    var queries = context.RequestServices.GetRequiredService<Queries>();
    var form = await context.Request.ReadFormAsync();
    var description = (form["description"].ToString() ?? string.Empty).Trim();
    var dueDateText = (form["due_date"].ToString() ?? string.Empty).Trim();

    var errors = new List<string>();
    if (description.Length == 0)
    {
        errors.Add("Description is required.");
    }
    if (description.Length > 500)
    {
        errors.Add("Description must be 500 characters or fewer.");
    }
    if (dueDateText.Length > 0)
    {
        var isValidFormat = Regex.IsMatch(dueDateText, @"^\d{4}-\d{2}-\d{2}$");
        var isValidDate = DateOnly.TryParseExact(dueDateText, "yyyy-MM-dd", CultureInfo.InvariantCulture, DateTimeStyles.None, out _);
        if (!isValidFormat || !isValidDate)
        {
            errors.Add("Due date must be a valid date (YYYY-MM-DD).");
        }
    }

    if (errors.Count > 0)
    {
        context.Response.Headers["HX-Retarget"] = "#form-errors";
        context.Response.Headers["HX-Reswap"] = "innerHTML";
        return Results.Extensions.RazorSlice<CSharpTodoApp.Slices.FormErrors, FormErrorsModel>(
            new FormErrorsModel(errors)
        );
    }

    await queries.CreateTodoAsync(description, "todo", dueDateText.Length > 0 ? dueDateText : null);
    var todos = await queries.ListTodosAsync(false, "", "");

    context.Response.Headers["HX-Trigger"] = "clearErrors";
    return Results.Extensions.RazorSlice<CSharpTodoApp.Slices.TodoList, TodoListModel>(
        new TodoListModel(todos, "", "")
    );
}));

app.MapPatch("/todos/{id:long}/toggle", async (long id, Queries queries) =>
{
    var todo = await queries.ToggleTodoStatusAsync(id);
    return todo is null
        ? Results.NotFound("Not Found")
        : Results.Extensions.RazorSlice<CSharpTodoApp.Slices.TodoItem, TodoItemModel>(
            new TodoItemModel(todo)
        );
});

app.MapDelete("/todos/{id:long}", async (long id, Queries queries) =>
{
    await queries.DeleteTodoAsync(id);
    return Results.Ok();
});

app.Run();

static bool IsHtmxRequest(HttpRequest request)
{
    return string.Equals(request.Headers["HX-Request"], "true", StringComparison.OrdinalIgnoreCase);
}

static void InitializeDatabase(string dataDir, string dbPath, string contentRootPath)
{
    Directory.CreateDirectory(dataDir);

    using var connection = new Microsoft.Data.Sqlite.SqliteConnection($"Data Source={dbPath}");
    connection.Open();

    using (var pragma = connection.CreateCommand())
    {
        pragma.CommandText = "PRAGMA journal_mode = WAL; PRAGMA foreign_keys = ON;";
        pragma.ExecuteNonQuery();
    }

    var ddl = File.ReadAllText(Path.Combine(contentRootPath, "schema.sql"));
    using var schemaCommand = connection.CreateCommand();
    schemaCommand.CommandText = ddl;
    schemaCommand.ExecuteNonQuery();
}
