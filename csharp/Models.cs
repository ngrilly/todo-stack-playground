namespace CSharpTodoApp;

public record Todo(
    long Id,
    string Description,
    string Status,
    string? DueDate,
    string CreatedAt
);

public record TodoPageModel(
    IReadOnlyList<Todo> Todos,
    string Filter,
    string Sort,
    IReadOnlyList<string> FormErrors
);

public record TodoListModel(
    IReadOnlyList<Todo> Todos,
    string Filter,
    string Sort
);

public record TodoItemModel(Todo Todo);

public record FormErrorsModel(IReadOnlyList<string> Errors);
