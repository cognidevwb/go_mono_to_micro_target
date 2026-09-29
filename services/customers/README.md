# customers-service

The customers service owns the customer record: it registers new customers (`POST /api/customers`) and looks them up by id (`GET /api/customers/:id`). It also answers whether a customer is active (`Service.IsActive`), the check other contexts such as orders rely on before accepting an order. It owns the `Customer` table and calls no other service.

| | |
|---|---|
| Owns tables | Customer |
| Aggregates | Customer |
| Calls | — |
| Orchestrates | — |
| Routes | /api/customers |
| Carved from | internal/customers/handlers.go, internal/customers/model.go, internal/customers/service.go |

## Run

```sh
cd services/customers && go run ./cmd/customers
```

Configuration is environment-only (`internal/config`); see the workspace `.env.example`.
