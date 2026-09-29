# inventory-service

The inventory service owns stock levels: each `StockItem` tracks the units on hand and the units reserved for a product. It reserves quantity for an order only when `on_hand - reserved` covers it, and releases the reservation through a compensating command if the order saga fails. It also reacts to the `order.placed` event and runs a periodic restock job that tops up items whose reserved count exceeds what is on hand.

| | |
|---|---|
| Owns tables | StockItem |
| Aggregates | StockItem |
| Calls | — |
| Orchestrates | — |
| Routes | /api/inventory |
| Carved from | internal/inventory/model.go, internal/inventory/service.go |

## Run

```sh
cd services/inventory && go run ./cmd/inventory
```

Configuration is environment-only (`internal/config`); see the workspace `.env.example`.
