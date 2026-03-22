from htpy import (
    Node,
    body,
    button,
    div,
    form,
    fragment,
    h1,
    head,
    html,
    input,
    label,
    li,
    link,
    main,
    meta,
    p,
    path,
    script,
    span,
    svg,
    title,
    ul,
)

from queries import Todo


def layout(page_title: str, content: Node) -> Node:
    return html(lang="en")[
        head[
            meta(charset="UTF-8"),
            meta(name="viewport", content="width=device-width, initial-scale=1.0"),
            title[page_title],
            script(src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"),
            link(
                rel="stylesheet",
                href="https://cdn.jsdelivr.net/npm/basecoat-css@0.3.11/dist/basecoat.cdn.min.css",
            ),
            script(
                src="https://cdn.jsdelivr.net/npm/basecoat-css@0.3.11/dist/js/all.min.js",
                defer=True,
            ),
            script(src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js"),
        ],
        body[main(".max-w-2xl.mx-auto.px-4.py-8")[content]],
    ]


def _filter_btn_class(active: bool) -> str:
    return "btn-sm" if active else "btn-sm-outline"


def todo_page(
    todos: list[Todo], filter: str, sort: str, form_errors: list[str]
) -> Node:
    return layout(
        "Todo App",
        [
            h1(".text-2xl.font-bold.mb-6")["Todo List"],
            div("#form-errors")[form_errors_view(form_errors)],
            todo_form(),
            div("#todo-list")[todo_list(todos, filter, sort)],
        ],
    )


def form_errors_view(errors: list[str]) -> Node:
    if not errors:
        return None
    return div(".alert.alert-destructive.mb-4", role="alert")[
        ul(".list-disc.list-inside")[[li[err] for err in errors]]
    ]


def todo_form() -> Node:
    return form(
        hx_post="/todos",
        hx_target="#todo-list",
        hx_swap="innerHTML",
        class_="card p-4 mb-4",
    )[
        div(".flex.items-end.gap-3")[
            div(".flex-1.flex.flex-col.gap-1\\.5")[
                label(".label", for_="description")["Description"],
                input(
                    ".input.w-full",
                    type="text",
                    id="description",
                    name="description",
                    placeholder="What needs to be done?",
                    required=True,
                    minlength="1",
                    maxlength="500",
                ),
            ],
            div(".flex.flex-col.gap-1\\.5")[
                label(".label", for_="due_date")["Due Date"],
                input(".input", type="date", id="due_date", name="due_date"),
            ],
            button(".btn", type="submit")["Add"],
        ]
    ]


def filter_bar(current_filter: str, current_sort: str) -> Node:
    return div(".flex.flex-wrap.items-center.gap-4.mb-4")[
        span(".text-sm.font-medium.text-muted-foreground")["Filter:"],
        div(role="group", class_="button-group")[
            button(
                class_=_filter_btn_class(current_filter == ""),
                hx_get=f"/?filter=&sort={current_sort}",
                hx_target="#todo-list",
                hx_swap="innerHTML",
            )["All"],
            button(
                class_=_filter_btn_class(current_filter == "todo"),
                hx_get=f"/?filter=todo&sort={current_sort}",
                hx_target="#todo-list",
                hx_swap="innerHTML",
            )["Todo"],
            button(
                class_=_filter_btn_class(current_filter == "done"),
                hx_get=f"/?filter=done&sort={current_sort}",
                hx_target="#todo-list",
                hx_swap="innerHTML",
            )["Done"],
        ],
        span(".text-sm.font-medium.text-muted-foreground")["Sort:"],
        div(role="group", class_="button-group")[
            button(
                class_=_filter_btn_class(current_sort == "description"),
                hx_get=f"/?filter={current_filter}&sort=description",
                hx_target="#todo-list",
                hx_swap="innerHTML",
            )["Description"],
            button(
                class_=_filter_btn_class(current_sort == "due_date"),
                hx_get=f"/?filter={current_filter}&sort=due_date",
                hx_target="#todo-list",
                hx_swap="innerHTML",
            )["Due Date"],
        ],
    ]


def todo_list(todos: list[Todo], filter: str, sort: str) -> Node:
    return fragment[
        filter_bar(filter, sort),
        (
            p(".text-muted-foreground.text-center.py-8")["No todos found."]
            if not todos
            else div(".flex.flex-col.gap-1")[
                [todo_item(todo) for todo in todos]
            ]
        ),
    ]


def todo_item(todo: Todo) -> Node:
    description_class = (
        ".flex-1.min-w-0.truncate.line-through.text-muted-foreground"
        if todo.status == "done"
        else ".flex-1.min-w-0.truncate"
    )
    return div(
        ".flex.items-center.gap-3.px-2.py-1\\.5.rounded-md.hover\\:bg-muted\\/50",
        id=f"todo-{todo.id}",
    )[
        input(
            type="checkbox",
            class_="checkbox",
            checked=todo.status == "done",
            hx_patch=f"/todos/{todo.id}/toggle",
            hx_target=f"#todo-{todo.id}",
            hx_swap="outerHTML",
        ),
        span(description_class)[todo.description],
        todo.due_date
        and span(".text-sm.text-muted-foreground.whitespace-nowrap")[todo.due_date],
        button(
            class_="btn-sm-icon-destructive",
            hx_delete=f"/todos/{todo.id}",
            hx_target=f"#todo-{todo.id}",
            hx_swap="outerHTML",
            hx_confirm="Delete this todo?",
        )[
            svg(
                xmlns="http://www.w3.org/2000/svg",
                width="24",
                height="24",
                viewBox="0 0 24 24",
                fill="none",
                stroke="currentColor",
                stroke_width="2",
                stroke_linecap="round",
                stroke_linejoin="round",
            )[
                path(d="M3 6h18"),
                path(d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"),
                path(d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"),
            ]
        ],
    ]
