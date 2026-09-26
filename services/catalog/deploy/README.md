# Deploying catalog-service

- Image: `docker build -f services/catalog/Dockerfile .` (from the workspace root).
- Kubernetes: `deploy/helm/catalog-service` (probes `/healthz`, `/readyz`; non-root, read-only root filesystem).
- Azure Container Apps: the `catalog-service` module in `deploy/azure/main.bicep` (`azd up`).
- Configuration: environment only — `DATABASE_URL`, `BROKER_URL`, `OTEL_EXPORTER_OTLP_ENDPOINT`.
