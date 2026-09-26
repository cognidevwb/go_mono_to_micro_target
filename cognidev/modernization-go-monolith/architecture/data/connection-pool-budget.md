---
id: connection-pool-budget
title: Connection-pool budget
family: data
currency: standard
role: invariant
applies_when:
  - answer:orchestrator != compose
teaches: Pods × pool size must fit under the database’s max_connections — set it, don’t inherit a default.
---

# Connection-pool budget

## What it is

Each replica opens its own pool. The sum across services and replica ceilings must stay under Postgres `max_connections` (default 100). `database/sql` defaults to unlimited open connections; pgxpool to max(4, NumCPU).

## Currency (2026)

**STANDARD — invariant** when autoscaled.

## Fit signals (from `.cognidev/feature/context.md`)

- `db_max_connections`, `load_profile` answers; replica ceilings in the HPA.

## What it changes in the generated services

- `DB_MAX_CONNS` env per service computed from the budget; `SetMaxOpenConns`/`MaxConns` set explicitly.

## Go 2026 implementation

```go
cfg, _ := pgxpool.ParseConfig(url); cfg.MaxConns = int32(maxConns)
// gorm: sqlDB, _ := db.DB(); sqlDB.SetMaxOpenConns(maxConns); sqlDB.SetConnMaxIdleTime(5*time.Minute)
```

## Anti-patterns

- Unlimited `database/sql` pools behind an autoscaler.

## Interacts with

- [[autoscaling]] · [[database-per-service]]
