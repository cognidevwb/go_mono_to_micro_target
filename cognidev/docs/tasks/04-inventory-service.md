# 4. Inventory Service

Extracted at position 4 of 5. Owns 1 table; makes 0 outbound calls.

## Owns

- `StockItem` — aggregate root

## Sagas this service is in

- participates in **CreateOrder**

## Moved from

- `internal/inventory/model.go`
- `internal/inventory/service.go`

This is a sample of what the map assigned; the porter moves everything
it owns, not just these.

## Files this task writes

| File | Layer | Purpose |
| --- | --- | --- |
| `services/inventory/internal/inventory/model.go` | data | model.go ported VERBATIM from the monolith's internal/inventory/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (inventory-service) |
| `services/inventory/internal/inventory/service.go` | services | service.go ported VERBATIM from the monolith's internal/inventory/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (inventory-service) |
| `services/inventory/cmd/inventory/main.go` | routing | inventory's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (inventory-service) |
| `services/inventory/migrations/0001_init.sql` | data | CREATE TABLE for the tables inventory OWNS and nothing else: stock_items (StockItem). Replaces the monolith's shared AutoMigrate; a column another service needs is served by inventory's API, never by a cross-schema join (inventory-service) |
| `services/inventory/internal/store/db.go` | data | open inventory's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY StockItem; the monolith's shared platform.Open and global AutoMigrate do not come over (inventory-service) |
| `services/inventory/internal/httpapi/stock_item_handlers.go` | routing | HTTP handlers for StockItem on the monolith's own router, kept (gin) under /v1/stockitem: bind + validate the request → call the ported service → JSON response, errors as RFC 9457 problem+json. context.Context flows from the request into every call. Routes carry the /v1 segment. (inventory-service) |
| `services/inventory/internal/events/order_placed_consumer.go` | services | consume `order.placed` (raised by orders) over NATS JetStream (github.com/nats-io/nats.go/jetstream: a durable consumer per subscriber, Nats-Msg-Id header for publish de-duplication, explicit Ack after the local commit): an idempotent handler keyed on the message id (processed-message table), acking only after the local commit — replaces the in-process bus Subscribe (inventory-service) |
| `services/inventory/internal/jobs/scheduler.go` | services | the scheduled work inventory inherits (StartRestockJob): run it exactly once across replicas — a single-replica worker, a Postgres advisory-lock lease, or a Kubernetes CronJob — started from cmd/inventory with the service's context so SIGTERM stops it (inventory-service) |
| `services/inventory/internal/acl/legacy.go` | services | anti-corruption layer for inventory: translate the monolith's shapes for StockItem into inventory's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (inventory-service) |
| `services/inventory/internal/inventory/inventory_test.go` | services | table-driven tests for inventory (StockItem): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (inventory-service) |
| `services/inventory/deploy/deploy.targets.yml` | cicd | CI/CD: fill inventory's deploy environments + post-deploy smoke-check path (# TODO(cognidev)) targeting Azure Container Apps (deploy/azure — azd), images in ACR; build (go build, go vet, go test -race), the distroless image, the Syft SBOM and the Grype scan are already wired (inventory-service) |
| `services/inventory/deploy/strangler.values.yaml` | deploy | Deployment: set inventory's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in inventory's coupling — a leaf higher, the hub / saga orchestrator canary low (inventory-service) |
| `services/inventory/README.md` | docs | Docs: fill inventory's <!-- CW-SEAM[kind=business-capability service=inventory] --> — 2-3 sentences on the business capability inventory owns, grounded in its source (inventory-service) |
| `services/inventory/test/equivalence/inventory_equivalence_test.go` | equivalence | Equivalence inventory: for each endpoint derive a {request, seed, expectedResponse} golden from the monolith's handler into testdata/equivalence/, replay it against the service with httptest, normalise volatile fields, assert new == legacy (inventory-service) |

