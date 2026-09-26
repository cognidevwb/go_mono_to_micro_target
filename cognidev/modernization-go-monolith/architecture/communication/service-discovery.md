---
id: service-discovery
title: Service discovery
family: communication
currency: declining
role: default
applies_when:
  - always
teaches: Let the platform resolve names — compose service names, Kubernetes Services, Container Apps app names — and pass base URLs by environment.
---

# Service discovery

## What it is

How a client finds its upstream. In 2026 this is a platform concern: DNS names provided by compose, Kubernetes Services or Azure Container Apps internal ingress. The application only reads a base URL from its environment.

## Currency (2026)

**DECLINING as an application concern** — client-side registries (Consul SDKs, Eureka) have given way to platform DNS.

## Fit signals (from `.cognidev/feature/context.md`)

- Every client edge.

## What it changes in the generated services

- `<CTX>_URL` env per upstream in compose, Helm values and ACA app settings, all using the same logical name (`http://catalog-service:8080`).

## Go 2026 implementation

```go
cfg.CatalogURL = envOr("CATALOG_URL", "http://catalog-service:8080")
```

## Anti-patterns

- Hardcoded hostnames.
- A service registry library in every service when the platform already resolves DNS.

## Interacts with

- [[sync-transport]] · [[deployment]]
