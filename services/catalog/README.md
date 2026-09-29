# catalog-service

The catalog service owns the product and category data: it lists products, creates new ones (rejecting non-positive prices), and serves `/api/products` for both. It also answers price lookups for other contexts through `PriceOf`, backed by a per-replica price cache whose entries expire after 30 seconds so replicas converge.

| | |
|---|---|
| Owns tables | Category, Product |
| Aggregates | Category |
| Calls | — |
| Orchestrates | — |
| Routes | /api/products |
| Carved from | internal/catalog/handlers.go, internal/catalog/model.go, internal/catalog/service.go |

## Run

```sh
cd services/catalog && go run ./cmd/catalog
```

Configuration is environment-only (`internal/config`); see the workspace `.env.example`.
