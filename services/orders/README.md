# orders-service

orders-service owns the `Order` and `OrderLine` records and exposes order placement (`POST /api/orders`) and a customer's own order history (`GET /api/orders/mine`). Placing an order looks up the customer and product prices through the customers and catalog clients, then runs the `CreateOrder` saga: stock is reserved in inventory first and payment is charged last, with compensations run in reverse on failure. The order and its lines are written only in the final local commit, together with an outbox row, so a failed saga leaves no order behind.

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
