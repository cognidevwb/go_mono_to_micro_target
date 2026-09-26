# Runbook

## Start locally

```sh
cp .env.example .env
docker compose up --build
```

Gateway: http://localhost:8080 · Grafana: http://localhost:3000

## Services

- `catalog-service` — `services/catalog`
- `customers-service` — `services/customers`
- `inventory-service` — `services/inventory`
- `payments-service` — `services/payments`
- `orders-service` — `services/orders`

Every service serves `/healthz` (process) and `/readyz` (its database). A `/readyz` 503 names the failed dependency.

## Rollback

Repoint the gateway's prefix at `LEGACY_UPSTREAM` (strangler) — the monolith still owns the route until its service is cut over.
