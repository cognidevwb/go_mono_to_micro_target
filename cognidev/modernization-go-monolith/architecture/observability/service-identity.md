---
id: service-identity
title: Service-to-service identity
family: observability
currency: rising
role: choice
applies_when:
  - services >= 6
teaches: At fleet scale, prove which service is calling — mesh mTLS or workload identity tokens, not network location.
---

# Service-to-service identity

## What it is

Each workload has a cryptographic identity (SPIFFE via mesh mTLS, or cloud workload identity tokens) and services authorize callers by it.

## Currency (2026)

**RISING, situational.**

## Fit signals (from `.cognidev/feature/context.md`)

- Six or more services; zero-trust requirements.

## What it changes in the generated services

- Mesh-provided mTLS (Istio ambient / Linkerd) or `github.com/spiffe/go-spiffe/v2` in clients.

## Go 2026 implementation

`source, _ := workloadapi.NewX509Source(ctx); tlsCfg := tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeID(catalogID))`

## Anti-patterns

- Hand-rolled shared secrets between services.

## Interacts with

- [[service-mesh]] · [[security]]
