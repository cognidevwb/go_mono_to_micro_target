---
id: deployment
title: Deployment — compose, Kubernetes, Azure Container Apps
family: reliability
currency: standard
role: default
applies_when:
  - always
teaches: One static binary per service in a distroless nonroot image, deployed the same way everywhere — compose locally, Kubernetes/Helm or Container Apps in production.
---

# Deployment — compose, Kubernetes, Azure Container Apps

## What it is

Each service builds with `CGO_ENABLED=0` into `gcr.io/distroless/static-debian12:nonroot`; compose runs the fleet locally with Postgres + NATS + collector; Kubernetes manifests/Helm or Azure Container Apps (Bicep/azd) in production.

## Currency (2026)

**STANDARD.**

## Fit signals (from `.cognidev/feature/context.md`)

- `orchestrator`, `cloud`, `deploy_packaging` answers.

## What it changes in the generated services

- Per-service Dockerfile; `deploy/compose.yaml`; `deploy/k8s/` or `deploy/helm/`; `deploy/azure/` for ACA.

## Go 2026 implementation

```dockerfile
FROM golang:1.26 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/svc ./services/orders/cmd/orders
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/svc /svc
USER nonroot:nonroot
ENTRYPOINT ["/svc"]
```

## Anti-patterns

- `:latest` tags; root users; shells in production images.

## Interacts with

- [[autoscaling]] · [[progressive-delivery]] · [[health-checks]]
