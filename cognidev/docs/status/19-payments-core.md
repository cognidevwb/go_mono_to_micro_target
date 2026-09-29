# Task 19 — Implement payments-service — domain + wiring

**COMPLETED**

Edited services/payments/internal/acl/legacy.go (8 files); compiled in the end-of-run pass.

## Files this task owned

- `services/payments/internal/acl/legacy.go` — Payments Service
- `services/payments/internal/app/create_order_compensation.go` — Payments Service
- `services/payments/internal/app/remote_operations.go` — Payments Service
- `services/payments/internal/httpapi/payment_handlers.go` — Payments Service
- `services/payments/internal/payments/payments_test.go` — Payments Service
- `services/payments/internal/store/db.go` — Payments Service
- `services/payments/migrations/0001_init.sql` — Payments Service
- `services/payments/test/equivalence/payments_equivalence_test.go` — Payments Service

## What it was asked to do

Resolve every remaining marker in this service's ported code: event seams publish through internal/events.Publish (the outbox) and consume from the broker; shared-state seams move to a shared store or are accepted explicitly; wiring seams port what the service needs from the monolith's main. Fill any `TODO(cognidev)` stub. Keep the ported code otherwise as it is — the monolith's own servers and routes that wire.go mounts stay mounted; never replace them with new handlers. Build and vet stay green.

## What it had to satisfy

- No CW-SEAM or TODO(cognidev) marker remains in these files.
- go build ./... and go vet ./... stay green.

---

Planned in `../tasks/`. Recorded from the engine's own per-task outcome;
nothing here was re-derived by reading the tree.
