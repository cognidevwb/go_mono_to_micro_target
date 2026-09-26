# catalog-service

The catalog service owns product and category data for the shop: it stores each product's SKU, name, price and category assignment, and each category's name and products. It exposes `GET /api/products` to list the catalog and `POST /api/products` to create a new product with basic validation (required SKU/name, positive price). It has no outbound calls to other services — it is a leaf context that other services (e.g. orders) read via the catalog client.

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
