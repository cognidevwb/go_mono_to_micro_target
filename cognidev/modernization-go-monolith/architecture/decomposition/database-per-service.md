---
id: database-per-service
title: Database per service
family: decomposition
currency: standard
role: invariant
applies_when:
  - always
teaches: Each service owns its tables; nobody else reads them — the one invariant that makes independent deployment real.
---

# Database per service

## What it is

Every service owns its schema and is the only process that connects to it. Other services get the data through its API or its events, never through a JOIN or a second `*gorm.DB` pointed at the same database. In the monolith one `platform.Open()` migrated every model into one database — that single handle is the seam.

## Currency (2026)

**STANDARD — invariant.** microservices.io, Azure Architecture Center. Logical databases on one Postgres server are an accepted first step; separate credentials are not optional.

## Fit signals (from `.cognidev/feature/context.md`)

- One `*gorm.DB` / `*sql.DB` shared by every package (`AutoMigrate(&Product{}, &Order{}, ...)`).
- Cross-context `Preload`/`Joins` or raw SQL joining two contexts' tables.

## What it changes in the generated services

- Each service: its own `DATABASE_URL`, user, and `migrations/0001_init.sql` containing ONLY its owned tables.
- Foreign keys across contexts are dropped; the ID stays as a plain column (`OrderLine.ProductID uint`).
- `plan.json.services[].entities` is the ownership list the migration generator reads.

## Go 2026 implementation

```go
// services/orders/cmd/orders/main.go
db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
// or: pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
```
Migrations: `golang-migrate` or `goose` SQL files per service, applied at startup or as a Job.

## Anti-patterns

- Two services with credentials to one schema "temporarily".
- `AutoMigrate` of another context's struct to make a query work.
- Sharing the monolith's migration history — each service starts its own.

## Interacts with

- [[saga-orchestration]] · [[api-composition]] · [[connection-pool-budget]] · [[expand-contract]]
