# Task 5 — Implement customers-service — domain + wiring

**COMPLETED**

Edited services/customers/internal/acl/legacy.go (6 files); compiled in the end-of-run pass.

## Files this task owned

- `services/customers/internal/acl/legacy.go` — Customers Service
- `services/customers/internal/customers/customers_test.go` — Customers Service
- `services/customers/internal/store/db.go` — Customers Service
- `services/customers/migrations/0001_init.sql` — Customers Service
- `services/customers/test/equivalence/customers_equivalence_test.go` — Customers Service

## What it was asked to do

Resolve every remaining marker in this service's ported code: event seams publish through internal/events.Publish (the outbox) and consume from the broker; shared-state seams move to a shared store or are accepted explicitly; wiring seams port what the service needs from the monolith's main. Fill any `TODO(cognidev)` stub. Keep the ported code otherwise as it is — the monolith's own servers and routes that wire.go mounts stay mounted; never replace them with new handlers. Build and vet stay green.

## What it had to satisfy

- No CW-SEAM or TODO(cognidev) marker remains in these files.
- go build ./... and go vet ./... stay green.

---

Planned in `../tasks/`. Recorded from the engine's own per-task outcome;
nothing here was re-derived by reading the tree.
