# Create the catalog slice — 8 files (data)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: model.go, 0001_init.sql, db.go, service.go, handlers.go, main.go, legacy.go, catalog_test.go. Read the grounding FIRST — the monolith sources this is ported from + the structural shards — and write every file against the others' real signatures (they are generated together, so make them consistent). Port the real behaviour from the monolith; do not invent APIs. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- model.go: model.go ported VERBATIM from the monolith's internal/catalog/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (catalog-service)
- 0001_init.sql: CREATE TABLE for the tables catalog OWNS and nothing else: categories (Category), products (Product). Replaces the monolith's shared AutoMigrate; a column another service needs is served by catalog's API, never by a cross-schema join (catalog-service)
- db.go: open catalog's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY Category, Product; the monolith's shared platform.Open and global AutoMigrate do not come over (catalog-service)
- service.go: service.go ported VERBATIM from the monolith's internal/catalog/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (catalog-service)
- handlers.go: handlers.go ported VERBATIM from the monolith's internal/catalog/handlers.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (catalog-service)
- main.go: catalog's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (catalog-service)
- legacy.go: anti-corruption layer for catalog: translate the monolith's shapes for Category, Product into catalog's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (catalog-service)
- catalog_test.go: table-driven tests for catalog (Category): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (catalog-service)

## Restricted access — edit only these (focus set)
- `services/catalog/internal/catalog/model.go`
- `services/catalog/migrations/0001_init.sql`
- `services/catalog/internal/store/db.go`
- `services/catalog/internal/catalog/service.go`
- `services/catalog/internal/catalog/handlers.go`
- `services/catalog/cmd/catalog/main.go`
- `services/catalog/internal/acl/legacy.go`
- `services/catalog/internal/catalog/catalog_test.go`

_Expansion: none — do not edit existing files in this step._

## Acceptance criteria (done when)
- All 8 files in the slice are written coherently against each other's real signatures.
- Ported from the monolith sources in the grounding — no invented APIs, reuses the app's ORM/auth libraries.
- The whole module compiles clean (the verification gate is green) after this task.

## Read first (grounding)
- .cognidev/understand/graph.json
- .cognidev/understand/state.json
- .cognidev/understand/flows.json
- .cognidev/understand/uml.json
- .cognidev/understand/files.json
- internal/catalog/model.go
- internal/catalog/service.go
- internal/catalog/handlers.go
