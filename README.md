# Platform

Microservices carved from the monolith — one Go module per service in a `go.work` workspace, behind a reverse-proxy gateway.

## Services

| Service | Module | Gateway prefix | Owns | Calls |
|---|---|---|---|---|
| `catalog-service` | [`services/catalog`](services/catalog/README.md) | `/api/products` | Category, Product | — |
| `customers-service` | [`services/customers`](services/customers/README.md) | `/api/customers` | Customer | — |
| `inventory-service` | [`services/inventory`](services/inventory/README.md) | `/api/inventory` | StockItem | — |
| `payments-service` | [`services/payments`](services/payments/README.md) | `/api/payments` | Payment | — |
| `orders-service` | [`services/orders`](services/orders/README.md) | `/api/orders` | Order, OrderLine | catalog, customers |

## Sagas

- `CreateOrder` — orchestrated by `orders`, participants: inventory, payments

## Build and run

```sh
make sync build vet test      # every module on its own
cp .env.example .env && docker compose up --build
```

See ARCHITECTURE.md, STATE-FLOW.md, RUNBOOK.md and docs/adr/.

The monolith's original README is kept at [docs/monolith/README.md](docs/monolith/README.md).
