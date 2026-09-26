# customers-service

The customers service owns customer identity and registration: it registers new customers with a unique, validated email and a display name, and serves lookups by customer ID. Each `Customer` record tracks an `Active` flag and creation timestamp, giving the rest of the platform a single source of truth for who a customer is without exposing the underlying storage.

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
