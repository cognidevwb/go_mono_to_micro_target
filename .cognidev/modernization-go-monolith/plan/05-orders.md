# Create the orders slice — 14 files (data)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: model.go, 0001_init.sql, db.go, outbox.go, service.go, create_order.go, create_order_test.go, handlers.go, main.go, catalog_client.go, customers_client.go, order_placed.go, legacy.go, orders_test.go. Read the grounding FIRST — the monolith sources this is ported from + the structural shards — and write every file against the others' real signatures (they are generated together, so make them consistent). Port the real behaviour from the monolith; do not invent APIs. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- model.go: model.go ported VERBATIM from the monolith's internal/orders/model.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (orders-service)
- 0001_init.sql: CREATE TABLE for the tables orders OWNS and nothing else: orders (Order), order_lines (OrderLine). Replaces the monolith's shared AutoMigrate; a column another service needs is served by orders's API, never by a cross-schema join (orders-service)
- db.go: open orders's OWN *gorm.DB (gorm.io/driver/postgres) from its own DSN — models registered here are ONLY Order, OrderLine; the monolith's shared platform.Open and global AutoMigrate do not come over (orders-service)
- outbox.go: transactional outbox for orders: an outbox table written in the SAME transaction as the business change, and a relay goroutine that publishes unsent rows over NATS JetStream (github.com/nats-io/nats.go/jetstream: a durable consumer per subscriber, Nats-Msg-Id header for publish de-duplication, explicit Ack after the local commit) with the row id as the de-duplication id, marking them sent only after the broker acks (orders-service)
- service.go: service.go ported VERBATIM from the monolith's internal/orders/service.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (orders-service) ALSO, in this file — the one change beyond its markers: replace the cross-service unit of work in PlaceOrder with the CreateOrder saga (outbox + compensations) — writes in inventory, orders, payments. Callers see what the monolith's single transaction showed them: the monolith rolled a failed order back, so it never existed — when the saga fails, after its compensations ran, delete the records this call created (the order and its lines) in this service's own transaction, so no list or get ever returns them; each compensation calls the participant's compensating endpoint, never only a log line.
- create_order.go: the CreateOrder ORCHESTRATED saga replacing the monolith's single db.Transaction across inventory → orders → payments: persist saga state in orders's database, send each step as a command over NATS JetStream (github.com/nats-io/nats.go/jetstream: a durable consumer per subscriber, Nats-Msg-Id header for publish de-duplication, explicit Ack after the local commit), the hardest-to-reverse step (the payment charge) LAST as the pivot, compensations in reverse (release stock, void payment) — each participant's Compensate calls that participant's POST /v1/<ctx>/create_order/compensate through a method you add to its typed client under internal/clients/<ctx>; a Compensate that only logs leaves the step committed. A step that makes one call per line item undoes the items it already did before returning its error — the monolith's transaction rolled those back too. State change + outgoing message commit together through the outbox. A failed saga leaves no business row the monolith's rollback would not have left: the rows the transaction created (the order and its lines) are written only in the final local commit, after the pivot succeeds — never saved first as pending and then marked cancelled or failed, which every read and list would show. Saga progress lives in the saga's own state table, not in the business row's status. Run(ctx, cmd) executes the steps in order and returns the final Outcome (Completed or Compensated) — the test calls Run and asserts the Outcome and the persisted saga state; it never waits for a published event. // CW-SEAM[kind=saga] (orders-service)
- create_order_test.go: CreateOrder saga test with fake participants: the happy path completes across inventory → orders → payments; a failure at the pivot runs the compensations in reverse and leaves no order row behind. Run(ctx, cmd) executes the steps in order and returns the final Outcome (Completed or Compensated) — the test calls Run and asserts the Outcome and the persisted saga state; it never waits for a published event (orders-service)
- handlers.go: handlers.go ported VERBATIM from the monolith's internal/orders/handlers.go by port-monolith-go — package clause kept, import paths rewritten to this module. Do not rewrite it: fill only its `// CW-SEAM[...]` markers (orders-service)
- main.go: orders's entrypoint: slog JSON handler, OpenTelemetry (OTLP) setup, its OWN database handle, the router, /healthz + /readyz, graceful shutdown on SIGTERM. Wire only this service's packages — the monolith's composition root is not ported (orders-service)
- catalog_client.go: typed client to catalog-service: net/http with otelhttp.NewTransport, a per-call context timeout, bounded retries on idempotent GETs only, base URL from config (service discovery), decoding this service's own copy of the contract in internal/contracts/catalog. Replaces the in-process call into the monolith's catalog package. A cross-context WRITE is a saga step, not a call. // CW-SEAM[kind=client-call] (orders-service)
- customers_client.go: typed client to customers-service: net/http with otelhttp.NewTransport, a per-call context timeout, bounded retries on idempotent GETs only, base URL from config (service discovery), decoding this service's own copy of the contract in internal/contracts/customers. Replaces the in-process call into the monolith's customers package. A cross-context WRITE is a saga step, not a call. // CW-SEAM[kind=client-call] (orders-service)
- order_placed.go: publish `order.placed` over NATS JetStream (github.com/nats-io/nats.go/jetstream: a durable consumer per subscriber, Nats-Msg-Id header for publish de-duplication, explicit Ack after the local commit) — written to this service's outbox table in the SAME transaction as the change it announces, relayed afterwards; replaces the in-process bus Publish. Payload is a versioned struct in this service's internal/events package, wrapped in the pkg/events envelope (orders-service)
- legacy.go: anti-corruption layer for orders: translate the monolith's shapes for Order, OrderLine into orders's model — stringly-typed statuses to typed constants, legacy column names, invariants rejected at the boundary. Pure functions, no I/O. // CW-SEAM[kind=acl-translate] (orders-service)
- orders_test.go: table-driven tests for orders (Order): the happy path per aggregate, validation → 400, not-found → 404, and the concurrency / out-of-stock guard; database via testcontainers-go Postgres, skipped with t.Skip when no container runtime is present (orders-service)

## Restricted access — edit only these (focus set)
- `services/orders/internal/orders/model.go`
- `services/orders/migrations/0001_init.sql`
- `services/orders/internal/store/db.go`
- `services/orders/internal/outbox/outbox.go`
- `services/orders/internal/orders/service.go`
- `services/orders/internal/saga/create_order.go`
- `services/orders/internal/saga/create_order_test.go`
- `services/orders/internal/orders/handlers.go`
- `services/orders/cmd/orders/main.go`
- `services/orders/internal/clients/catalog_client.go`
- `services/orders/internal/clients/customers_client.go`
- `services/orders/internal/events/order_placed.go`
- `services/orders/internal/acl/legacy.go`
- `services/orders/internal/orders/orders_test.go`

_Expansion: none — do not edit existing files in this step._

## Acceptance criteria (done when)
- All 14 files in the slice are written coherently against each other's real signatures.
- Ported from the monolith sources in the grounding — no invented APIs, reuses the app's ORM/auth libraries.
- The whole module compiles clean (the verification gate is green) after this task.

## Read first (grounding)
- .cognidev/understand/graph.json
- .cognidev/understand/state.json
- .cognidev/understand/flows.json
- .cognidev/understand/uml.json
- .cognidev/understand/files.json
- internal/orders/model.go
- internal/orders/service.go
- internal/orders/handlers.go
- services/orders/internal/contracts/catalog/contracts.go
- services/orders/internal/contracts/customers/contracts.go
