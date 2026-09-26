---
id: saga-isolation
title: Saga isolation countermeasures
family: data
currency: standard
role: invariant
applies_when:
  - sagas > 0
teaches: A saga has no isolation — other requests see its intermediate states, so design semantic locks and reversible steps.
---

# Saga isolation countermeasures

## What it is

Between saga steps other transactions see partial state: stock reserved for an order that will fail. Countermeasures: semantic locks (`status = pending`), commutative updates, reordering (pivot step last), and re-reading values at the step that needs them. The fixture's `Reserve` does read-check-write — two sagas can oversell; the fix is an atomic conditional update.

## Currency (2026)

**STANDARD, under-applied.** Garcia-Molina & Salem; Richardson *Microservices Patterns* ch. 4.

## Fit signals (from `.cognidev/feature/context.md`)

- Any saga; read-then-write sequences in a saga step (`entity-access.json`).

## What it changes in the generated services

- Pending states on the orchestrator's aggregate; conditional updates in participants.

## Go 2026 implementation

```go
tag, err := tx.Exec(ctx, `UPDATE stock_items SET reserved = reserved + $2 WHERE product_id = $1 AND on_hand - reserved >= $2`, id, qty)
if tag.RowsAffected() == 0 { return ErrOutOfStock }
```

## Anti-patterns

- Assuming the saga is serializable because each step is.
- Compensation that overwrites a concurrent change.

## Interacts with

- [[saga-orchestration]] · [[optimistic-concurrency]]
