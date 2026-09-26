# Architecture

A Go 1.26 `go.work` workspace: one module per service, a shared kernel (`pkg/`)
kept to transport concerns, and a go-proxy gateway at the edge. Router: gin. Data: gorm on
PostgreSQL, one database per service. Events: nats through a transactional outbox.

```mermaid
flowchart LR
  client([client]) --> gateway
  gateway --> catalog[catalog-service]
  catalog --> catalogdb[(catalogdb)]
  gateway --> customers[customers-service]
  customers --> customersdb[(customersdb)]
  gateway --> inventory[inventory-service]
  inventory --> inventorydb[(inventorydb)]
  gateway --> payments[payments-service]
  payments --> paymentsdb[(paymentsdb)]
  gateway --> orders[orders-service]
  orders --> ordersdb[(ordersdb)]
  orders -. http .-> catalog
  orders -. http .-> customers
  broker{{nats}}
  catalog <-- events --> broker
  customers <-- events --> broker
  inventory <-- events --> broker
  payments <-- events --> broker
  orders <-- events --> broker
```

| Module | Owns tables | Calls | Routes |
|---|---|---|---|
| `services/catalog` | Category, Product | — | /api/products |
| `services/customers` | Customer | — | /api/customers |
| `services/inventory` | StockItem | — | /api/inventory |
| `services/payments` | Payment | — | /api/payments |
| `services/orders` | Order, OrderLine | catalog, customers | /api/orders |

## Sagas

- **CreateOrder** — orchestrated by `orders`: inventory → payments

## Why these boundaries

The monolith's five internal packages (`catalog`, `customers`, `inventory`, `payments`,
`orders`) mapped one-to-one onto the five services above — none were merged or split.
Each package already owned a disjoint set of gorm models (`service-map.json`) and had no
package-level state shared with another, so the module boundary could follow the existing
package boundary exactly.

`catalog`, `customers`, `inventory` and `payments` had no outbound calls into a sibling
package in the monolith, so they became leaf services: each owns its own table(s) and
exposes a REST API, with no service-to-service coupling to redesign.

`orders` was the one context coupled to others — it called into `catalog` (to price and
validate order lines) and `customers` (to validate the customer) synchronously, and its
checkout flow wrote to `orders`, `inventory` and `payments` inside a single database
transaction. That transaction could not be preserved once the tables moved to separate
databases, so it was redesigned into the `CreateOrder` saga orchestrated by `orders`:
inventory is reserved first (reversible), payment is captured last (hardest to reverse),
and a failure compensates by releasing the reservation. The synchronous reads into
`catalog` and `customers` were kept as accepted coupling — they are simple, low-latency
lookups with no write side effects — and became HTTP calls through generated clients
(`internal/clients/catalog`, `internal/clients/customers`) rather than events, since
`orders` needs an authoritative, synchronous answer before accepting an order.
