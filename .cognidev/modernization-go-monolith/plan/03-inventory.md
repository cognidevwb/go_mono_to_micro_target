# Create the inventory slice — 10 files (data)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: model.go, 0001_init.sql, db.go, service.go, order_placed_consumer.go, scheduler.go, legacy.go, main.go, stock_item_handlers.go, inventory_test.go. Read the grounding FIRST — the monolith sources this is ported from + the structural shards — and write every file against the others' real signatures (they are generated together, so make them consistent). Port the real behaviour from the monolith; do not invent APIs. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- model.go: model.go ported VERBATIM from the monolith's internal/inventory/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (inventory-service)
- 0001_init.sql: CREATE TABLE for the tables inventory OWNS and nothing else: stock_items (StockItem). Replaces the monolith's shared AutoMigrate; a column another service needs is served by inventory's API, never by a cross-schema join (inventory-service)
- db.go: open inventory's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY StockItem; the monolith's shared platform.Open and global AutoMigrate do not come over (inventory-service)
- service.go: service.go ported VERBATIM from the monolith's internal/inventory/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (inventory-service)
- order_placed_consumer.go: consume `order.placed` (raised by orders) over NATS JetStream (github.com/nats-io/nats.go/jetstream: a durable consumer per subscriber, Nats-Msg-Id header for publish de-duplication, explicit Ack after the local commit): an idempotent handler keyed on the message id (processed-message table), acking only after the local commit — replaces the in-process bus Subscribe (inventory-service)
- scheduler.go: the scheduled work inventory inherits (StartRestockJob): run it exactly once across replicas — a single-replica worker, a Postgres advisory-lock lease, or a Kubernetes CronJob — started from cmd/inventory with the service's context so SIGTERM stops it (inventory-service)
- legacy.go: anti-corruption layer for inventory: translate the monolith's shapes for StockItem into inventory's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (inventory-service)
- main.go: inventory's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (inventory-service)
- stock_item_handlers.go: HTTP handlers for StockItem on the monolith's own router, kept (gin) under /v1/stockitem: bind + validate the request → call the ported service → JSON response, errors as RFC 9457 problem+json. context.Context flows from the request into every call. Routes carry the /v1 segment. (inventory-service)
- inventory_test.go: table-driven tests for inventory (StockItem): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (inventory-service)

## Restricted access — edit only these (focus set)
- `services/inventory/internal/inventory/model.go`
- `services/inventory/migrations/0001_init.sql`
- `services/inventory/internal/store/db.go`
- `services/inventory/internal/inventory/service.go`
- `services/inventory/internal/events/order_placed_consumer.go`
- `services/inventory/internal/jobs/scheduler.go`
- `services/inventory/internal/acl/legacy.go`
- `services/inventory/cmd/inventory/main.go`
- `services/inventory/internal/httpapi/stock_item_handlers.go`
- `services/inventory/internal/inventory/inventory_test.go`

_Expansion: none — do not edit existing files in this step._

## Acceptance criteria (done when)
- All 10 files in the slice are written coherently against each other's real signatures.
- Ported from the monolith sources in the grounding — no invented APIs, reuses the app's ORM/auth libraries.
- The whole module compiles clean (the verification gate is green) after this task.

## Read first (grounding)
- .cognidev/understand/graph.json
- .cognidev/understand/state.json
- .cognidev/understand/flows.json
- .cognidev/understand/uml.json
- .cognidev/understand/files.json
- internal/inventory/model.go
- internal/inventory/service.go
