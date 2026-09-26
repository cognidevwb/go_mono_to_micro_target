# Task 9 — Implement inventory-service — domain + wiring

**COMPLETED**

Edited services/inventory/internal/acl/legacy.go (10 files); compiled in the end-of-run pass.

## Files this task owned

- `services/inventory/internal/acl/legacy.go` — Inventory Service
- `services/inventory/internal/app/wire.go` — Inventory Service
- `services/inventory/internal/events/order_placed_consumer.go` — Inventory Service
- `services/inventory/internal/httpapi/stock_item_handlers.go` — Inventory Service
- `services/inventory/internal/inventory/inventory_test.go` — Inventory Service
- `services/inventory/internal/jobs/scheduler.go` — Inventory Service
- `services/inventory/internal/store/db.go` — Inventory Service
- `services/inventory/migrations/0001_init.sql` — Inventory Service
- `services/inventory/test/equivalence/inventory_equivalence_test.go` — Inventory Service

## What it was asked to do

Resolve every remaining marker in this service's ported code: event seams publish through internal/events.Publish (the outbox) and consume from the broker; shared-state seams move to a shared store or are accepted explicitly; wiring seams port what the service needs from the monolith's main. Fill any `TODO(cognidev)` stub. Keep the ported code otherwise as it is — the monolith's own servers and routes that wire.go mounts stay mounted; never replace them with new handlers. Build and vet stay green.

## What it had to satisfy

- No CW-SEAM or TODO(cognidev) marker remains in these files.
- go build ./... and go vet ./... stay green.

---

Planned in `../tasks/`. Recorded from the engine's own per-task outcome;
nothing here was re-derived by reading the tree.
