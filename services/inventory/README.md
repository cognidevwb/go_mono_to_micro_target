# inventory-service

The inventory service owns on-hand and reserved stock counts for each product (`StockItem`). It reserves units against incoming orders, restocks items whose on-hand quantity has fallen below what is reserved, and consumes `order.placed` events from the orders context to keep stock levels in sync with order activity.

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
