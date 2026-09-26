# Deploying payments-service

- Image: `docker build -f services/payments/Dockerfile .` (from the workspace root).
- Kubernetes: `deploy/helm/payments-service` (probes `/healthz`, `/readyz`; non-root, read-only root filesystem).
- Azure Container Apps: the `payments-service` module in `deploy/azure/main.bicep` (`azd up`).
- Configuration: environment only — `DATABASE_URL`, `BROKER_URL`, `OTEL_EXPORTER_OTLP_ENDPOINT`.
