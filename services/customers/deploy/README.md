# Deploying customers-service

- Image: `docker build -f services/customers/Dockerfile .` (from the workspace root).
- Kubernetes: `deploy/helm/customers-service` (probes `/healthz`, `/readyz`; non-root, read-only root filesystem).
- Azure Container Apps: the `customers-service` module in `deploy/azure/main.bicep` (`azd up`).
- Configuration: environment only — `DATABASE_URL`, `BROKER_URL`, `OTEL_EXPORTER_OTLP_ENDPOINT`.
