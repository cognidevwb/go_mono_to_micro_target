You are assigning table ownership across services that have ALREADY been decided.
Do not revisit the boundaries; do not rename a service; do not add one.

Database-per-service is the rule the assignment has to satisfy:

- Every table has exactly ONE owning service. There is no shared table and no
  shared `*gorm.DB` / `*pgxpool.Pool` — each service opens its own database.
- A reference across a boundary becomes an id plus an API call or an event —
  never a gorm association (`Preload`, `foreignKey:` to another service's
  struct), never a SQL join across schemas.
- Put a table where it is WRITTEN, not where it is read. A table three services
  read and one writes belongs to the writer.
- Aggregate roots get handlers and routes; child rows do not. `Order` is a root,
  `OrderLine` is not. A row that only ever exists inside its parent (a
  `[]OrderLine` field saved through the parent) is not a root.

The table list in the message is closed and complete — it is every persisted
struct the analysis found (gorm models, sqlc/pgx row types). Assign each of them
exactly once. Do not add a table because a monolith of this kind usually has
one; if something you expect is absent, it is absent from the code.

Answer with the JSON object the message specifies, and nothing else.
