# Agent Instructions

## sqlc: Dynamic Queries with Boolean Flags and CASE Expressions

sqlc does not support dynamic SQL (dynamic WHERE clauses, ORDER BY, etc.).
To work around this, use boolean flags with CASE expressions as described
in [Sqlc: 2024 check in](https://brandur.org/fragments/sqlc-2024).

See `go/db/query.sql` (`ListTodos`) for a working example of dynamic
filtering and sorting using this technique. Note the SQLite-specific
workaround: `sqlc.arg()` is not expanded in ORDER BY clauses, so the
sort parameter is registered via a dummy WHERE condition and referenced
by positional parameter (`?N`) in ORDER BY.
