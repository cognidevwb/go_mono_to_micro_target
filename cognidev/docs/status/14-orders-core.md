# Task 14 — Implement orders-service — domain + wiring

**COMPLETED**

Edited services/orders/internal/acl/legacy.go — this task's work is complete. 1 seam(s) still in it belong to 1 later task(s) in this plan: 15-orders-saga (1).

## Files this task owned

- `services/orders/internal/acl/legacy.go` — Orders Service
- `services/orders/internal/events/order_placed.go` — Orders Service
- `services/orders/internal/orders/orders_test.go` — Orders Service
- `services/orders/internal/orders/service.go` — Orders Service
- `services/orders/internal/outbox/outbox.go` — Orders Service
- `services/orders/internal/store/db.go` — Orders Service
- `services/orders/migrations/0001_init.sql` — Orders Service
- `services/orders/test/equivalence/orders_equivalence_test.go` — Orders Service

## What it was asked to do

Resolve every remaining marker in this service's ported code: event seams publish through internal/events.Publish (the outbox) and consume from the broker; shared-state seams move to a shared store or are accepted explicitly; wiring seams port what the service needs from the monolith's main. Fill any `TODO(cognidev)` stub. Keep the ported code otherwise as it is — the monolith's own servers and routes that wire.go mounts stay mounted; never replace them with new handlers. Build and vet stay green.

## What it had to satisfy

- No CW-SEAM or TODO(cognidev) marker remains in these files.
- go build ./... and go vet ./... stay green.

---

Planned in `../tasks/`. Recorded from the engine's own per-task outcome;
nothing here was re-derived by reading the tree.
