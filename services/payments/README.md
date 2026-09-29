# payments-service

The payments service captures money for an order and records the result as a `Payment` row (order id, amount, status `captured`). It wraps the payment gateway behind `Service.Charge`, which only persists the payment once the gateway charge succeeds. It is the last, hardest-to-reverse step of the order saga, so it also exposes a compensating operation for a failed order.

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
