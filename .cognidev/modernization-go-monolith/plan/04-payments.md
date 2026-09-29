# Create the payments slice — 12 files (data)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: model.go, 0001_init.sql, db.go, gateway.go, service.go, legacy.go, create_order_compensation.go, main.go, payment_handlers.go, wire.go, remote_operations.go, payments_test.go. Read the grounding FIRST — the monolith sources this is ported from + the structural shards — and write every file against the others' real signatures (they are generated together, so make them consistent). Port the real behaviour from the monolith; do not invent APIs. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- model.go: model.go ported VERBATIM from the monolith's internal/payments/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (payments-service)
- 0001_init.sql: CREATE TABLE for the tables payments OWNS and nothing else: payments (Payment). Replaces the monolith's shared AutoMigrate; a column another service needs is served by payments's API, never by a cross-schema join (payments-service)
- db.go: open payments's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY Payment; the monolith's shared platform.Open and global AutoMigrate do not come over (payments-service)
- gateway.go: gateway.go ported VERBATIM from the monolith's internal/payments/gateway.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (payments-service)
- service.go: service.go ported VERBATIM from the monolith's internal/payments/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (payments-service)
- legacy.go: anti-corruption layer for payments: translate the monolith's shapes for Payment into payments's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (payments-service)
- create_order_compensation.go: the compensating operation payments offers the CreateOrder saga: undo exactly what payments's step did (a stock reservation is released, a charge is voided) in payments's own transaction, a no-op when repeated. Served as POST /v1/payments/create_order/compensate, mounted from this file's init func by appending to the package's `mounts` (routes.go) — wire.go is not edited. A record it keeps to stay idempotent (a ledger of undone steps) is a `CREATE TABLE IF NOT EXISTS` appended to the package's `schemas` from init — app.go creates it at startup; a table nothing creates fails the first compensation. The orchestrator's Compensate for payments calls it through its typed client — a log line is not a compensation (payments-service)
- main.go: payments's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (payments-service)
- payment_handlers.go: HTTP handlers for Payment on the monolith's own router, kept (gin) under /v1/payment: bind + validate the request → call the ported service → JSON response, errors as RFC 9457 problem+json. context.Context flows from the request into every call. The package exports `RegisterRoutes(r gin.IRouter, db *gorm.DB, svc …)`, which internal/app/wire.go calls. Routes carry the /v1 segment. (payments-service)
- wire.go: mount the new HTTP handlers: the handlers package (internal/httpapi) exports ONE entry point, `RegisterRoutes(r gin.IRouter, db *gorm.DB, svc <the ported service type>)`, and wire() calls it once after the auth middleware, with the service value wire already constructs. Keep every other line of wire as port-monolith-go wrote it. Without this call no route of payments is served (payments-service)
- remote_operations.go: serve the ported methods other services call on payments: `Charge` (called by orders, a write inside a saga) at POST /v1/payments/charge. For each: bind the arguments as JSON, call the ported method unchanged, answer its results as JSON; an error answers {"error": err.Error()} with the status the monolith's own handlers use for it, so the caller sees the monolith's message. Mounted from this file's init by appending to the package's `mounts` (routes.go) — wire.go is not edited (payments-service)
- payments_test.go: table-driven tests for payments (Payment): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (payments-service)

## Restricted access — edit only these (focus set)
- `services/payments/internal/payments/model.go`
- `services/payments/migrations/0001_init.sql`
- `services/payments/internal/store/db.go`
- `services/payments/internal/payments/gateway.go`
- `services/payments/internal/payments/service.go`
- `services/payments/internal/acl/legacy.go`
- `services/payments/internal/app/create_order_compensation.go`
- `services/payments/cmd/payments/main.go`
- `services/payments/internal/httpapi/payment_handlers.go`
- `services/payments/internal/app/wire.go`
- `services/payments/internal/app/remote_operations.go`
- `services/payments/internal/payments/payments_test.go`

_Expansion: none — do not edit existing files in this step._

## Acceptance criteria (done when)
- All 12 files in the slice are written coherently against each other's real signatures.
- Ported from the monolith sources in the grounding — no invented APIs, reuses the app's ORM/auth libraries.
- The whole module compiles clean (the verification gate is green) after this task.

## Read first (grounding)
- .cognidev/understand/graph.json
- .cognidev/understand/state.json
- .cognidev/understand/flows.json
- .cognidev/understand/uml.json
- .cognidev/understand/files.json
- internal/payments/model.go
- internal/payments/gateway.go
- internal/payments/service.go
