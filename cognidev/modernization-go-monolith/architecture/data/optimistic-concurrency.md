---
id: optimistic-concurrency
title: Optimistic concurrency
family: data
currency: standard
role: invariant
applies_when:
  - always
teaches: A version column and a conditional UPDATE — the split widened every write window, and last-writer-wins now loses data.
---

# Optimistic concurrency

## What it is

Each aggregate row carries a `version`; updates are `WHERE id = $1 AND version = $2` and increment it; zero rows affected means a concurrent change — reload and retry or return 409.

## Currency (2026)

**STANDARD — invariant.**

## Fit signals (from `.cognidev/feature/context.md`)

- Read-modify-write sequences (`First` then `Save` in gorm) in `entity-access.json`.

## What it changes in the generated services

- `version int` on owned aggregates; `If-Match`/ETag on write routes.

## Go 2026 implementation

```go
res := tx.Model(&o).Where("version = ?", o.Version).Updates(map[string]any{"status": s, "version": gorm.Expr("version + 1")})
if res.RowsAffected == 0 { return ErrConflict }
```

## Anti-patterns

- gorm `Save` of a whole struct read earlier (silently overwrites).

## Interacts with

- [[saga-isolation]] · [[expand-contract]]
