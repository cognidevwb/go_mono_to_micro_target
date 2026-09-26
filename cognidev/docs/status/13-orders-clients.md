# Task 13 — Implement orders-service — cross-service calls

**COMPLETED**

Edited services/orders/internal/clients/catalog/client.go — this task's work is complete. 17 seam(s) still in it belong to 3 later task(s) in this plan: 14-orders-core (13), 15-orders-saga (1), 19-payments-core (3).

## Files this task owned

- `services/orders/internal/clients/catalog/client.go` — Orders Service
- `services/orders/internal/clients/catalog_contract_test.go` — Orders Service
- `services/orders/internal/clients/customers/client.go` — Orders Service
- `services/orders/internal/clients/customers_contract_test.go` — Orders Service
- `services/orders/internal/clients/inventory/client.go` — Orders Service
- `services/orders/internal/clients/payments/client.go` — Orders Service

## What it was asked to do

Implement each `// CW-SEAM[kind=cross-call]` stub in these generated clients: a real HTTP call to the provider's endpoint (baseURL(), the httpx client), mapping the JSON response into the declared result types. Delete each marker. Every module still builds (`go build ./... && go vet ./...`).

## What it had to satisfy

- No CW-SEAM[kind=cross-call] marker remains in the clients.
- go build ./... and go vet ./... stay green.

---

Planned in `../tasks/`. Recorded from the engine's own per-task outcome;
nothing here was re-derived by reading the tree.
