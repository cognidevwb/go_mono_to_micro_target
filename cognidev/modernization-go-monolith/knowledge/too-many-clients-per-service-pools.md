---
id: too-many-clients-per-service-pools
runtime: go
since: 0
category: runtime-failure
tags: [postgres, pgx, gorm, connection-pool, autoscaling, database-per-service]
symptoms:
  - "sorry, too many clients already"
  - "SQLSTATE 53300"
  - "remaining connection slots are reserved"
  - "sql: database is closed"
remedy: { kind: environment }
reference: "PostgreSQL max_connections (default 100) · pgxpool Config.MaxConns default = max(4, runtime.NumCPU())"
---

# One monolith pool became N services × replicas pools

## The signature

```
failed to connect to `user=orders database=orders`: server error: FATAL: sorry, too many clients already (SQLSTATE 53300)
```

or, after a graceful shutdown closed the pool while a background job still ran:

```
restock failed: sql: database is closed
```

## What it means

The monolith opened ONE `*gorm.DB` (`database/sql` defaults to unlimited open
connections) against one Postgres. After the cut each service opens its own
pool, each autoscaled replica opens another, and when the logical databases
share one server (the common first step) they share its `max_connections`.
Five services × 4 replicas × 25 connections is 500 against a default of 100.

`sql: database is closed` is the other half: a goroutine the monolith started
(a ticker job) outlives the `db.Close()` in the service's shutdown path.

## What to do

- Budget explicitly: `SetMaxOpenConns` / `pgxpool MaxConns` = floor(limit ×
  share / max replicas), from the `db_max_connections` answer.
- Background jobs take the service's root `context.Context` and exit before the
  pool is closed (`errgroup` + `signal.NotifyContext`).
