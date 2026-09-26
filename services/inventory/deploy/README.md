# Deploying inventory-service

- Image: `docker build -f services/inventory/Dockerfile .` (from the workspace root).
- Kubernetes: `deploy/helm/inventory-service` (probes `/healthz`, `/readyz`; non-root, read-only root filesystem).
- Azure Container Apps: the `inventory-service` module in `deploy/azure/main.bicep` (`azd up`).
- Configuration: environment only — `DATABASE_URL`, `BROKER_URL`, `OTEL_EXPORTER_OTLP_ENDPOINT`.
