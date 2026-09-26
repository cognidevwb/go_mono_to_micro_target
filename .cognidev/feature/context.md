# Decomposition brief — Go monolith → microservices

Deterministically distilled from the structural shards. Refine names, merge
over-split contexts, resolve each coupling edge, design the sagas, set the order.

## Candidate bounded contexts (each owns its tables)

### catalog (catalog → `catalog-service`)
- owned tables: Category, Product
- entity relationships:
    Category →(many) Product [navigation]
    Product →(one) Category [foreign-key]
    Product →(one) Category [navigation]
- HTTP surface (2):
    GET /api/products
    POST /api/products
- files (3) by role:
    db (1):
      internal/catalog/service.go
    http (1):
      internal/catalog/handlers.go
    other (1):
      internal/catalog/model.go

### customers (customers → `customers-service`)
- owned tables: Customer
- HTTP surface (2):
    GET /api/customers/:id
    POST /api/customers
- files (3) by role:
    db (1):
      internal/customers/service.go
    http (1):
      internal/customers/handlers.go
    other (1):
      internal/customers/model.go

### inventory (inventory → `inventory-service`)
- owned tables: StockItem
- entity relationships:
    StockItem →(one) Product [foreign-key]
- files (2) by role:
    db (1):
      internal/inventory/service.go
    other (1):
      internal/inventory/model.go

  - scheduled work it inherits: StartRestockJob
### payments (payments → `payments-service`)
- owned tables: Payment
- entity relationships:
    Payment →(one) Order [foreign-key]
- files (3) by role:
    db (1):
      internal/payments/service.go
    other (2):
      internal/payments/gateway.go
      internal/payments/model.go

  - external calls it inherits: POST /v1/charges via g.client
### orders (orders → `orders-service`)
- owned tables: Order, OrderLine
- entity relationships:
    Order →(many) OrderLine [navigation]
    Order →(one) Customer [foreign-key]
    OrderLine →(one) Order [foreign-key]
    OrderLine →(one) Product [foreign-key]
- HTTP surface (2):
    GET /api/orders/mine
    POST /api/orders
- files (3) by role:
    db (1):
      internal/orders/service.go
    http (1):
      internal/orders/handlers.go
    other (1):
      internal/orders/model.go

  - request context it reads today, which callers must now send: user (1)
## Cross-context coupling matrix (A → B (n) resolved calls)

- orders → catalog (1)
- orders → customers (1)
- orders → inventory (1)
- orders → payments (1)

## Events that stop being method calls

- `order.placed` — raised in orders, handled in inventory: needs a broker topic and an outbox

## What the sections above cover

- ambient reads: 1 of 2 sit in a candidate service; the other 1 is in internal/platform (1)
- event handlers: 1 of 1 sit in a candidate service
- cache accesses: 2 of 2 sit in a candidate service
- scheduled jobs: 1 of 1 sit in a candidate service
- outbound calls: 1 of 1 sit in a candidate service

## Tables no service owns

_(none — every table belongs to exactly one candidate)_

## State that stops being shared

One process keeps one copy of a package-level variable. Two services keep two, and they drift — a `sync.Mutex` beside it guards one process, not the fleet.

- `catalog` — `catalog.priceCache: map[uint]float64`

## Cache keys straddling a cut

_(none — no cache key is touched from two candidate services)_

## Endpoints that need more than one service

_(none of the 6 endpoints measured reaches two candidate services)_

## Saga candidates (a unit of work spanning ≥3 contexts)

- **CreateOrder** — orchestrator `orders-service`, participants ["orders","customers","inventory","catalog","payments"]

