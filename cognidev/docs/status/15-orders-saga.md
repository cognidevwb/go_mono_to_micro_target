# Task 15 — Implement orders-service — orchestration saga

**COMPLETED**

Edited services/orders/internal/saga/create_order.go; compiled in the end-of-run pass.

## Files this task owned

- `services/orders/internal/saga/create_order.go` — Orders Service

## What it was asked to do

Bind each saga participant to its generated client: a Do and a Compensate per step, the payment step last, state persisted with the outbox in one transaction. Replace the `// CW-SEAM[kind=saga]` transaction in the ported code with a call to the saga. Delete the markers.

## What it had to satisfy

- The saga binds every participant with a compensation; no marker remains.
- go test ./... passes.

---

Planned in `../tasks/`. Recorded from the engine's own per-task outcome;
nothing here was re-derived by reading the tree.
