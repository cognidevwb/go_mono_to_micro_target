# 2. Catalog Service

Extracted at position 2 of 5. Owns 2 tables; makes 0 outbound calls.

## Owns

- `Category` — aggregate root
- `Product`

## Moved from

- `internal/catalog/handlers.go`
- `internal/catalog/model.go`
- `internal/catalog/service.go`

This is a sample of what the map assigned; the porter moves everything
it owns, not just these.

## Files this task writes

| File | Layer | Purpose |
| --- | --- | --- |
| `services/catalog/internal/catalog/handlers.go` | routing | handlers.go ported VERBATIM from the monolith's internal/catalog/handlers.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (catalog-service) |
| `services/catalog/internal/catalog/model.go` | data | model.go ported VERBATIM from the monolith's internal/catalog/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (catalog-service) |
| `services/catalog/internal/catalog/service.go` | services | service.go ported VERBATIM from the monolith's internal/catalog/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (catalog-service) |
| `services/catalog/cmd/catalog/main.go` | routing | catalog's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (catalog-service) |
| `services/catalog/migrations/0001_init.sql` | data | CREATE TABLE for the tables catalog OWNS and nothing else: categories (Category), products (Product). Replaces the monolith's shared AutoMigrate; a column another service needs is served by catalog's API, never by a cross-schema join (catalog-service) |
| `services/catalog/internal/store/db.go` | data | open catalog's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY Category, Product; the monolith's shared platform.Open and global AutoMigrate do not come over (catalog-service) |
| `services/catalog/internal/acl/legacy.go` | services | anti-corruption layer for catalog: translate the monolith's shapes for Category, Product into catalog's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (catalog-service) |
| `services/catalog/internal/catalog/catalog_test.go` | services | table-driven tests for catalog (Category): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (catalog-service) |
| `services/catalog/deploy/deploy.targets.yml` | cicd | CI/CD: fill catalog's deploy environments + post-deploy smoke-check path (# TODO(cognidev)) targeting Azure Container Apps (deploy/azure — azd), images in ACR; build (go build, go vet, go test -race), the distroless image, the Syft SBOM and the Grype scan are already wired (catalog-service) |
| `services/catalog/deploy/strangler.values.yaml` | deploy | Deployment: set catalog's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in catalog's coupling — a leaf higher, the hub / saga orchestrator canary low (catalog-service) |
| `services/catalog/README.md` | docs | Docs: fill catalog's <!-- CW-SEAM[kind=business-capability service=catalog] --> — 2-3 sentences on the business capability catalog owns, grounded in its source (catalog-service) |
| `services/catalog/test/equivalence/catalog_equivalence_test.go` | equivalence | Equivalence catalog: for each endpoint derive a {request, seed, expectedResponse} golden from the monolith's handler into testdata/equivalence/, replay it against the service with httptest, normalise volatile fields, assert new == legacy (catalog-service) |

