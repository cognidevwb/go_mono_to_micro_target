# Deploying orders-service

- Image: `docker build -f services/orders/Dockerfile .` (from the workspace root).
- Kubernetes: `deploy/helm/orders-service` (probes `/healthz`, `/readyz`; non-root, read-only root filesystem).
- Azure Container Apps: the `orders-service` module in `deploy/azure/main.bicep` (`azd up`).
- Configuration: environment only — `DATABASE_URL`, `BROKER_URL`, `OTEL_EXPORTER_OTLP_ENDPOINT`, `CATALOG_SERVICE_URL`, `CUSTOMERS_SERVICE_URL`.
