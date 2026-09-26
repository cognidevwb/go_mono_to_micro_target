# payments-service

The payments service captures charges against orders: it calls out to the external payment provider through `internal/payments/gateway.go` (`Gateway.Charge`) and, on success, records the resulting `Payment` (order, amount, status) in its own table via `internal/payments/service.go` (`Service.Charge`). It owns the `Payment` aggregate exclusively and exposes it over `/api/payments`, acting as the hardest-to-reverse participant in the orders placement saga.

| | |
|---|---|
| Owns tables | Payment |
| Aggregates | Payment |
| Calls | — |
| Orchestrates | — |
| Routes | /api/payments |
| Carved from | internal/payments/gateway.go, internal/payments/model.go, internal/payments/service.go |

## Run

```sh
cd services/payments && go run ./cmd/payments
```

Configuration is environment-only (`internal/config`); see the workspace `.env.example`.
