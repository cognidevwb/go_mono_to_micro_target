# 1. Orders Service

Extracted at position 1 of 5. Owns 2 tables; makes 2 outbound calls.

## Owns

- `Order` — aggregate root
- `OrderLine`

## Talks to

- `catalog-service` — typed HttpClient with timeout, retry and circuit breaker
- `customers-service` — typed HttpClient with timeout, retry and circuit breaker

## Sagas this service is in

- orchestrates **CreateOrder**

## Moved from

- `internal/orders/handlers.go`
- `internal/orders/model.go`
- `internal/orders/service.go`

This is a sample of what the map assigned; the porter moves everything
it owns, not just these.

## Files this task writes

| File | Layer | Purpose |
| --- | --- | --- |
| `services/orders/internal/orders/handlers.go` | routing | handlers.go ported VERBATIM from the monolith's internal/orders/handlers.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (orders-service) |
| `services/orders/internal/orders/model.go` | data | model.go ported VERBATIM from the monolith's internal/orders/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (orders-service) |
| `services/orders/internal/orders/service.go` | services | service.go ported VERBATIM from the monolith's internal/orders/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (orders-service) |
| `services/orders/cmd/orders/main.go` | routing | orders's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (orders-service) |
| `services/orders/migrations/0001_init.sql` | data | CREATE TABLE for the tables orders OWNS and nothing else: orders (Order), order_lines (OrderLine). Replaces the monolith's shared AutoMigrate; a column another service needs is served by orders's API, never by a cross-schema join (orders-service) |
| `services/orders/internal/store/db.go` | data | open orders's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY Order, OrderLine; the monolith's shared platform.Open and global AutoMigrate do not come over (orders-service) |
| `services/orders/internal/clients/catalog_client.go` | routing | typed client to catalog-service: net/http with otelhttp.NewTransport, a per-call context timeout, bounded retries on idempotent GETs only, base URL from config (service discovery), decoding this service's own copy of the contract in internal/contracts/catalog. Replaces the in-process call into the monolith's catalog package. A cross-context WRITE is a saga step, not a call. // CW-SEAM[kind=client-call] (orders-service) |
| `services/orders/internal/clients/customers_client.go` | routing | typed client to customers-service: net/http with otelhttp.NewTransport, a per-call context timeout, bounded retries on idempotent GETs only, base URL from config (service discovery), decoding this service's own copy of the contract in internal/contracts/customers. Replaces the in-process call into the monolith's customers package. A cross-context WRITE is a saga step, not a call. // CW-SEAM[kind=client-call] (orders-service) |
| `services/orders/internal/events/order_placed.go` | services | publish `order.placed` over NATS JetStream (github.com/nats-io/nats.go/jetstream: a durable consumer per subscriber, Nats-Msg-Id header for publish de-duplication, explicit Ack after the local commit) — written to this service's outbox table in the SAME transaction as the change it announces, relayed afterwards; replaces the in-process bus Publish. Payload is a versioned struct in this service's internal/events package, wrapped in the pkg/events envelope (orders-service) |
| `services/orders/internal/acl/legacy.go` | services | anti-corruption layer for orders: translate the monolith's shapes for Order, OrderLine into orders's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (orders-service) |
| `services/orders/internal/orders/orders_test.go` | services | table-driven tests for orders (Order): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (orders-service) |
| `services/orders/internal/outbox/outbox.go` | data | transactional outbox for orders: an outbox table written in the SAME transaction as the business change, and a relay goroutine that publishes unsent rows over NATS JetStream (github.com/nats-io/nats.go/jetstream: a durable consumer per subscriber, Nats-Msg-Id header for publish de-duplication, explicit Ack after the local commit) with the row id as the de-duplication id, marking them sent only after the broker acks (orders-service) |
| `services/orders/internal/saga/create_order.go` | services | the CreateOrder ORCHESTRATED saga replacing the monolith's single db.Transaction across inventory → orders → payments: persist saga state in orders's database, send each step as a command over NATS JetStream (github.com/nats-io/nats.go/jetstream: a durable consumer per subscriber, Nats-Msg-Id header for publish de-duplication, explicit Ack after the local commit), the hardest-to-reverse step (the payment charge) LAST as the pivot, compensations in reverse (release stock, void payment). State change + outgoing message commit together through the outbox. Run(ctx, cmd) executes the steps in order and returns the final Outcome (Completed or Compensated) — the test calls Run and asserts the Outcome and the persisted saga state; it never waits for a published event. // CW-SEAM[kind=saga] (orders-service) |
| `services/orders/internal/saga/create_order_test.go` | services | CreateOrder saga test with fake participants: the happy path completes across inventory → orders → payments; a failure at the pivot runs the compensations in reverse. Run(ctx, cmd) executes the steps in order and returns the final Outcome (Completed or Compensated) — the test calls Run and asserts the Outcome and the persisted saga state; it never waits for a published event (orders-service) |
| `services/orders/deploy/deploy.targets.yml` | cicd | CI/CD: fill orders's deploy environments + post-deploy smoke-check path (# TODO(cognidev)) targeting Azure Container Apps (deploy/azure — azd), images in ACR; build (go build, go vet, go test -race), the distroless image, the Syft SBOM and the Grype scan are already wired (orders-service) |
| `services/orders/deploy/strangler.values.yaml` | deploy | Deployment: set orders's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in orders's coupling — a leaf higher, the hub / saga orchestrator canary low (orders-service) |
| `services/orders/README.md` | docs | Docs: fill orders's <!-- CW-SEAM[kind=business-capability service=orders] --> — 2-3 sentences on the business capability orders owns, grounded in its source (orders-service) |
| `services/orders/internal/clients/catalog_contract_test.go` | contracts | Contract tests orders->catalog: for each testdata/contracts/catalog/*.json golden response assert it decodes into internal/contracts/catalog with no unknown fields (json.Decoder.DisallowUnknownFields) and satisfies the value invariants; record one fixture (orders-service) |
| `services/orders/internal/clients/customers_contract_test.go` | contracts | Contract tests orders->customers: for each testdata/contracts/customers/*.json golden response assert it decodes into internal/contracts/customers with no unknown fields (json.Decoder.DisallowUnknownFields) and satisfies the value invariants; record one fixture (orders-service) |
| `services/orders/test/equivalence/orders_equivalence_test.go` | equivalence | Equivalence orders: for each endpoint derive a {request, seed, expectedResponse} golden from the monolith's handler into testdata/equivalence/, replay it against the service with httptest, normalise volatile fields, assert new == legacy (orders-service) |
| `services/orders/internal/contracts/catalog/contracts.go` | data | orders-service's own copy of the catalog service's contract: the request/response DTOs it decodes (JSON tags, no gorm tags, no behaviour), copied from the provider's API — never imported from another module |
| `services/orders/internal/contracts/customers/contracts.go` | data | orders-service's own copy of the customers service's contract: the request/response DTOs it decodes (JSON tags, no gorm tags, no behaviour), copied from the provider's API — never imported from another module |

