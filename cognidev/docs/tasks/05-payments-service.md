# 5. Payments Service

Extracted at position 5 of 5. Owns 1 table; makes 0 outbound calls.

## Owns

- `Payment` — aggregate root

## Sagas this service is in

- participates in **CreateOrder**

## Moved from

- `internal/payments/gateway.go`
- `internal/payments/model.go`
- `internal/payments/service.go`

This is a sample of what the map assigned; the porter moves everything
it owns, not just these.

## Files this task writes

| File | Layer | Purpose |
| --- | --- | --- |
| `services/payments/internal/payments/gateway.go` | services | gateway.go ported VERBATIM from the monolith's internal/payments/gateway.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (payments-service) |
| `services/payments/internal/payments/model.go` | data | model.go ported VERBATIM from the monolith's internal/payments/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (payments-service) |
| `services/payments/internal/payments/service.go` | services | service.go ported VERBATIM from the monolith's internal/payments/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (payments-service) |
| `services/payments/cmd/payments/main.go` | routing | payments's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (payments-service) |
| `services/payments/migrations/0001_init.sql` | data | CREATE TABLE for the tables payments OWNS and nothing else: payments (Payment). Replaces the monolith's shared AutoMigrate; a column another service needs is served by payments's API, never by a cross-schema join (payments-service) |
| `services/payments/internal/store/db.go` | data | open payments's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY Payment; the monolith's shared platform.Open and global AutoMigrate do not come over (payments-service) |
| `services/payments/internal/httpapi/payment_handlers.go` | routing | HTTP handlers for Payment on the monolith's own router, kept (gin) under /v1/payment: bind + validate the request → call the ported service → JSON response, errors as RFC 9457 problem+json. context.Context flows from the request into every call. The package exports `RegisterRoutes(r gin.IRouter, db *gorm.DB, svc …)`, which internal/app/wire.go calls. Routes carry the /v1 segment. (payments-service) |
| `services/payments/internal/app/wire.go` | routing | mount the new HTTP handlers: the handlers package (internal/httpapi) exports ONE entry point, `RegisterRoutes(r gin.IRouter, db *gorm.DB, svc <the ported service type>)`, and wire() calls it once after the auth middleware, with the service value wire already constructs. Keep every other line of wire as port-monolith-go wrote it. Without this call no route of payments is served (payments-service) |
| `services/payments/internal/app/remote_operations.go` | routing | serve the ported methods other services call on payments: `Charge` (called by orders, a write inside a saga) at POST /v1/payments/charge. For each: bind the arguments as JSON, call the ported method unchanged, answer its results as JSON; an error answers {"error": err.Error()} with the status the monolith's own handlers use for it, so the caller sees the monolith's message. Mounted from this file's init by appending to the package's `mounts` (routes.go) — wire.go is not edited (payments-service) |
| `services/payments/internal/acl/legacy.go` | services | anti-corruption layer for payments: translate the monolith's shapes for Payment into payments's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (payments-service) |
| `services/payments/internal/payments/payments_test.go` | services | table-driven tests for payments (Payment): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (payments-service) |
| `services/payments/internal/app/create_order_compensation.go` | services | the compensating operation payments offers the CreateOrder saga: undo exactly what payments's step did (a stock reservation is released, a charge is voided) in payments's own transaction, a no-op when repeated. Served as POST /v1/payments/create_order/compensate, mounted from this file's init func by appending to the package's `mounts` (routes.go) — wire.go is not edited. A record it keeps to stay idempotent (a ledger of undone steps) is a `CREATE TABLE IF NOT EXISTS` appended to the package's `schemas` from init — app.go creates it at startup; a table nothing creates fails the first compensation. The orchestrator's Compensate for payments calls it through its typed client — a log line is not a compensation (payments-service) |
| `services/payments/deploy/deploy.targets.yml` | cicd | CI/CD: fill payments's deploy environments + post-deploy smoke-check path (# TODO(cognidev)) targeting Azure Container Apps (deploy/azure — azd), images in ACR; build (go build, go vet, go test -race), the distroless image, the Syft SBOM and the Grype scan are already wired (payments-service) |
| `services/payments/deploy/strangler.values.yaml` | deploy | Deployment: set payments's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in payments's coupling — a leaf higher, the hub / saga orchestrator canary low (payments-service) |
| `services/payments/README.md` | docs | Docs: fill payments's <!-- CW-SEAM[kind=business-capability service=payments] --> — 2-3 sentences on the business capability payments owns, grounded in its source (payments-service) |
| `services/payments/test/equivalence/payments_equivalence_test.go` | equivalence | Equivalence payments: for each endpoint derive a {request, seed, expectedResponse} golden from the monolith's handler into testdata/equivalence/, replay it against the service with httptest, normalise volatile fields, assert new == legacy (payments-service) |

