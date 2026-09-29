# Create the customers slice — 9 files (data)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: model.go, 0001_init.sql, db.go, service.go, handlers.go, main.go, remote_operations.go, legacy.go, customers_test.go. Read the grounding FIRST — the monolith sources this is ported from + the structural shards — and write every file against the others' real signatures (they are generated together, so make them consistent). Port the real behaviour from the monolith; do not invent APIs. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- model.go: model.go ported VERBATIM from the monolith's internal/customers/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (customers-service)
- 0001_init.sql: CREATE TABLE for the tables customers OWNS and nothing else: customers (Customer). Replaces the monolith's shared AutoMigrate; a column another service needs is served by customers's API, never by a cross-schema join (customers-service)
- db.go: open customers's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY Customer; the monolith's shared platform.Open and global AutoMigrate do not come over (customers-service)
- service.go: service.go ported VERBATIM from the monolith's internal/customers/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (customers-service)
- handlers.go: handlers.go ported VERBATIM from the monolith's internal/customers/handlers.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (customers-service)
- main.go: customers's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (customers-service)
- remote_operations.go: serve the ported methods other services call on customers: `IsActive` (called by orders) at POST /v1/customers/is_active. For each: bind the arguments as JSON, call the ported method unchanged, answer its results as JSON; an error answers {"error": err.Error()} with the status the monolith's own handlers use for it, so the caller sees the monolith's message. Mounted from this file's init by appending to the package's `mounts` (routes.go) — wire.go is not edited (customers-service)
- legacy.go: anti-corruption layer for customers: translate the monolith's shapes for Customer into customers's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (customers-service)
- customers_test.go: table-driven tests for customers (Customer): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (customers-service)

## Restricted access — edit only these (focus set)
- `services/customers/internal/customers/model.go`
- `services/customers/migrations/0001_init.sql`
- `services/customers/internal/store/db.go`
- `services/customers/internal/customers/service.go`
- `services/customers/internal/customers/handlers.go`
- `services/customers/cmd/customers/main.go`
- `services/customers/internal/app/remote_operations.go`
- `services/customers/internal/acl/legacy.go`
- `services/customers/internal/customers/customers_test.go`

_Expansion: none — do not edit existing files in this step._

## Acceptance criteria (done when)
- All 9 files in the slice are written coherently against each other's real signatures.
- Ported from the monolith sources in the grounding — no invented APIs, reuses the app's ORM/auth libraries.
- The whole module compiles clean (the verification gate is green) after this task.

## Read first (grounding)
- .cognidev/understand/graph.json
- .cognidev/understand/state.json
- .cognidev/understand/flows.json
- .cognidev/understand/uml.json
- .cognidev/understand/files.json
- internal/customers/model.go
- internal/customers/service.go
- internal/customers/handlers.go
