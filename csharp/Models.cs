namespace CSharpTodoApp;

public sealed record Todo(
    long Id,
    string Description,
    string Status,
    string? DueDate,
    string CreatedAt
);

public sealed record TodoPageModel(
    IReadOnlyList<Todo> Todos,
    string Filter,
    string Sort,
    IReadOnlyList<string> FormErrors
);

public sealed record TodoListModel(
    IReadOnlyList<Todo> Todos,
    string Filter,
    string Sort
);

public sealed record TodoItemModel(Todo Todo);

public sealed record FormErrorsModel(IReadOnlyList<string> Errors);
