# 3. Customers Service

Extracted at position 3 of 5. Owns 1 table; makes 0 outbound calls.

## Owns

- `Customer` — aggregate root

## Moved from

- `internal/customers/handlers.go`
- `internal/customers/model.go`
- `internal/customers/service.go`

This is a sample of what the map assigned; the porter moves everything
it owns, not just these.

## Files this task writes

| File | Layer | Purpose |
| --- | --- | --- |
| `services/customers/internal/customers/handlers.go` | routing | handlers.go ported VERBATIM from the monolith's internal/customers/handlers.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (customers-service) |
| `services/customers/internal/customers/model.go` | data | model.go ported VERBATIM from the monolith's internal/customers/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (customers-service) |
| `services/customers/internal/customers/service.go` | services | service.go ported VERBATIM from the monolith's internal/customers/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (customers-service) |
| `services/customers/cmd/customers/main.go` | routing | customers's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (customers-service) |
| `services/customers/migrations/0001_init.sql` | data | CREATE TABLE for the tables customers OWNS and nothing else: customers (Customer). Replaces the monolith's shared AutoMigrate; a column another service needs is served by customers's API, never by a cross-schema join (customers-service) |
| `services/customers/internal/store/db.go` | data | open customers's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY Customer; the monolith's shared platform.Open and global AutoMigrate do not come over (customers-service) |
| `services/customers/internal/app/remote_operations.go` | routing | serve the ported methods other services call on customers: `IsActive` (called by orders) at POST /v1/customers/is_active. For each: bind the arguments as JSON, call the ported method unchanged, answer its results as JSON; an error answers {"error": err.Error()} with the status the monolith's own handlers use for it, so the caller sees the monolith's message. Mounted from this file's init by appending to the package's `mounts` (routes.go) — wire.go is not edited (customers-service) |
| `services/customers/internal/acl/legacy.go` | services | anti-corruption layer for customers: translate the monolith's shapes for Customer into customers's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (customers-service) |
| `services/customers/internal/customers/customers_test.go` | services | table-driven tests for customers (Customer): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (customers-service) |
| `services/customers/deploy/deploy.targets.yml` | cicd | CI/CD: fill customers's deploy environments + post-deploy smoke-check path (# TODO(cognidev)) targeting Azure Container Apps (deploy/azure — azd), images in ACR; build (go build, go vet, go test -race), the distroless image, the Syft SBOM and the Grype scan are already wired (customers-service) |
| `services/customers/deploy/strangler.values.yaml` | deploy | Deployment: set customers's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in customers's coupling — a leaf higher, the hub / saga orchestrator canary low (customers-service) |
| `services/customers/README.md` | docs | Docs: fill customers's <!-- CW-SEAM[kind=business-capability service=customers] --> — 2-3 sentences on the business capability customers owns, grounded in its source (customers-service) |
| `services/customers/test/equivalence/customers_equivalence_test.go` | equivalence | Equivalence customers: for each endpoint derive a {request, seed, expectedResponse} golden from the monolith's handler into testdata/equivalence/, replay it against the service with httptest, normalise volatile fields, assert new == legacy (customers-service) |

