# orders-service

Orders owns the customer's order placement and history: it validates that a customer is active, prices each requested line from catalog, and records the order and its lines in its own database. Because reserving stock and charging payment now live in separate services, placing an order runs the `CreateOrder` saga, which reserves inventory then charges payment — the harder-to-reverse step last — and compensates completed steps if a later one fails. Once the order is placed it is saved and an `order.placed` event is published through the transactional outbox for other services to consume.

| | |
|---|---|
| Owns tables | Order, OrderLine |
| Aggregates | Order |
| Calls | catalog, customers |
| Orchestrates | CreateOrder |
| Routes | /api/orders |
| Carved from | internal/orders/handlers.go, internal/orders/model.go, internal/orders/service.go |

## Run

```sh
cd services/orders && go run ./cmd/orders
```

Configuration is environment-only (`internal/config`); see the workspace `.env.example`.
